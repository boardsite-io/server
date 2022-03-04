package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/heat1q/boardsite/api/log"

	"github.com/heat1q/boardsite/session"

	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/websocket"
	"github.com/heat1q/boardsite/attachment"
	"github.com/labstack/echo/v4"
)

const (
	SessionCtxKey = "boardsite-session"
	SecretCtxKey  = "boardsite-session-secret"
	UserCtxKey    = "boardsite-user"
)

type Handler interface {
	PostCreateSession(c echo.Context) error
	PutSessionConfig(c echo.Context) error
	GetSessionConfig(c echo.Context) error
	PostUsers(c echo.Context) error
	PutUser(c echo.Context) error
	GetSocket(c echo.Context) error
	GetPageRank(c echo.Context) error
	PostPages(c echo.Context) error
	PutPages(c echo.Context) error
	GetPage(c echo.Context) error
	GetPageSync(c echo.Context) error
	PostPageSync(c echo.Context) error
	PostAttachment(c echo.Context) error
	GetAttachment(c echo.Context) error
}

type handler struct {
	cfg        config.Session
	dispatcher session.Dispatcher
}

func NewHandler(cfg config.Session, dispatcher session.Dispatcher) Handler {
	return &handler{
		cfg:        cfg,
		dispatcher: dispatcher,
	}
}

// PostCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
func (h *handler) PostCreateSession(c echo.Context) error {
	scb, err := h.dispatcher.Create(c.Request().Context(), session.NewConfig(h.cfg))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, session.CreateSessionResponse{
		Config: scb.Config(),
	})
}

func (h *handler) PutSessionConfig(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	var cfg session.Config
	if err := c.Bind(&cfg); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.SetConfig(cfg); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithErrorf("setConfig: %w", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetSessionConfig(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, session.GetConfigResponse{Users: scb.GetUsers(), Config: scb.Config()})
}

func (h *handler) PostUsers(c echo.Context) error {
	scb, err := h.dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	var userReq session.User

	if err := json.NewDecoder(c.Request().Body).Decode(&userReq); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	// new user struct with alias and color
	user, err := scb.NewUser(userReq.Alias, userReq.Color)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *handler) PutUser(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}
	if err := scb.KickUser(c.Param("userId")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetSocket handles request for a websocket upgrade
// based on the sessionID and the userID.
func (h *handler) GetSocket(c echo.Context) error {
	scb, err := h.dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}
	err = websocket.Subscribe(c, scb, c.Param("userId"))
	if err != nil {
		log.Ctx(c.Request().Context()).Errorf("websocket subscribe: %v", err)
	}
	return nil
}

func (h *handler) GetPageRank(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	pageRank, err := scb.GetPageRank(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pageRank)
}

// PostPages handles requests regarding adding or retrieving pages.
func (h *handler) PostPages(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	// add a Page
	var data session.PageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.AddPages(c.Request().Context(), data); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *handler) PutPages(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	op := c.QueryParam(session.QueryKeyUpdate)

	var data session.PageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.UpdatePages(c.Request().Context(), data, op); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetPage(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	pageID := c.Param("pageId")
	if !scb.IsValidPage(c.Request().Context(), pageID) {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	page, err := scb.GetPage(c.Request().Context(), pageID, true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, page)
}

func (h *handler) GetPageSync(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	pageRank, err := scb.GetPageRank(c.Request().Context())
	if err != nil {
		return err
	}

	sync, err := scb.GetPageSync(c.Request().Context(), pageRank, true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sync)
}

func (h *handler) PostPageSync(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	var sync session.PageSync
	if err := json.NewDecoder(c.Request().Body).Decode(&sync); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.SyncSession(c.Request().Context(), sync); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) PostAttachment(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	if err := c.Request().ParseMultipartForm(2 << 20); err != nil {
		return apiErrors.From(apiErrors.AttachmentSizeExceeded).Wrap(
			apiErrors.WithMessage("file size exceeded limit of 2MB"),
			apiErrors.WithError(err))
	}
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	attachID, err := scb.Attachments().Upload(data)
	if err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.JSON(http.StatusCreated, attachment.AttachmentResponse{AttachID: attachID})
}

func (h *handler) GetAttachment(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}
	attachID := c.Param("attachId")
	data, MIMEType, err := scb.Attachments().Get(attachID)
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	return c.Stream(http.StatusOK, MIMEType, data)
}

func getSCB(c echo.Context) (session.Controller, error) {
	scb, ok := c.Get(SessionCtxKey).(session.Controller)
	if !ok {
		return nil, apiErrors.ErrForbidden
	}
	return scb, nil
}
