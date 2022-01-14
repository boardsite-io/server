package session

import (
	"encoding/json"
	"io"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/redis"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	PostCreateSession(c echo.Context) error
	GetUsers(c echo.Context) error
	PostUsers(c echo.Context) error
	GetSocket(c echo.Context) error
	GetPages(c echo.Context) error
	PostPages(c echo.Context) error
	PutPages(c echo.Context) error
	DeletePages(c echo.Context) error
	GetPageUpdate(c echo.Context) error
	DeletePageUpdate(c echo.Context) error
	PostAttachment(c echo.Context) error
	GetAttachment(c echo.Context) error
}

type handler struct {
	cfg        *config.Configuration
	Dispatcher Dispatcher
}

func NewHandler(cfg *config.Configuration, cache redis.Handler) Handler {
	return &handler{
		cfg:        cfg,
		Dispatcher: NewDispatcher(cache),
	}
}

// PostCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
func (h *handler) PostCreateSession(c echo.Context) error {
	idstr, err := h.Dispatcher.Create(c.Request().Context(), h.cfg.Session.MaxUsers)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, types.CreateSessionResponse{SessionID: idstr})
}

func (h *handler) GetUsers(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}
	return c.JSON(http.StatusOK, scb.GetUsers())
}

func (h *handler) PostUsers(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	var userReq types.User

	if err := json.NewDecoder(c.Request().Body).Decode(&userReq); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	// new user struct with alias and color
	user, err := scb.NewUser(userReq.Alias, userReq.Color)
	if err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	return c.JSON(http.StatusCreated, user)
}

// GetSocket handles request for a websocket upgrade
// based on the sessionID and the userID.
func (h *handler) GetSocket(c echo.Context) error {
	var (
		sessionID = c.Param("id")
		userID    = c.Param("userId")
	)

	scb, err := h.Dispatcher.GetSCB(sessionID)
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}
	_, errUser := scb.GetUserReady(userID)
	if errUser != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	onConnect := func(conn *gws.Conn) error {
		return Subscribe(c.Request().Context(), conn, scb, userID)
	}

	return upgrade(c, onConnect)
}

func (h *handler) GetPages(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	pageRank, meta, err := scb.GetPages(c.Request().Context())
	if err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	// return pagerank array
	pages := types.ContentPageSync{
		PageRank: pageRank,
		Meta:     meta,
	}
	return c.JSON(http.StatusOK, pages)
}

// PostPages handles requests regarding adding or retrieving pages.
func (h *handler) PostPages(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	// add a Page
	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.AddPages(c.Request().Context(), data.PageID, data.Index, data.Meta); err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.NoContent(http.StatusCreated)
}

func (h *handler) PutPages(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.UpdatePages(c.Request().Context(), data.PageID, data.Meta, data.Clear); err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeletePages(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.ErrBadRequest.Wrap(apiErrors.WithError(err))
	}

	if err := scb.DeletePages(c.Request().Context(), data.PageID...); err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) GetPageUpdate(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	pageID := c.Param("pageId")
	if !scb.IsValidPage(c.Request().Context(), pageID) {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	strokes, errFetch := scb.GetStrokes(c.Request().Context(), pageID)
	if errFetch != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.JSON(http.StatusOK, strokes)
}

func (h *handler) DeletePageUpdate(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	pageID := c.Param("pageId")
	if !scb.IsValidPage(c.Request().Context(), pageID) {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	if err := scb.DeletePages(c.Request().Context(), pageID); err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) PostAttachment(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	if err := c.Request().ParseMultipartForm(2 << 20); err != nil {
		return apiErrors.ErrBadRequest.Wrap(
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

	attachID, err := scb.Attachments.Upload(data)
	if err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.JSON(http.StatusCreated, types.AttachmentResponse{AttachID: attachID})
}

func (h *handler) GetAttachment(c echo.Context) error {
	scb, err := h.Dispatcher.GetSCB(c.Param("id"))
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}
	attachID := c.Param("attachId")
	data, MIMEType, err := scb.Attachments.Get(attachID)
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	return c.Stream(http.StatusOK, MIMEType, data)
}
