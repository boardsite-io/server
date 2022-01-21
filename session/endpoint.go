package session

import (
	"encoding/json"
	"io"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/heat1q/boardsite/api/config"
	apiErrors "github.com/heat1q/boardsite/api/errors"
	"github.com/heat1q/boardsite/api/types"
	"github.com/labstack/echo/v4"
)

const (
	SessionCtxKey = "boardsite-session"
	UserCtxKey    = "boardsite-user"
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
	dispatcher Dispatcher
}

func NewHandler(cfg *config.Configuration, dispatcher Dispatcher) Handler {
	return &handler{
		cfg:        cfg,
		dispatcher: dispatcher,
	}
}

// PostCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
func (h *handler) PostCreateSession(c echo.Context) error {
	idstr, err := h.dispatcher.Create(c.Request().Context(), h.cfg.Session.MaxUsers)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, types.CreateSessionResponse{SessionID: idstr})
}

func (h *handler) PostUsers(c echo.Context) error {
	scb, err := h.dispatcher.GetSCB(c.Param("id"))
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
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *handler) GetUsers(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, scb.GetUsers())
}

// GetSocket handles request for a websocket upgrade
// based on the sessionID and the userID.
func (h *handler) GetSocket(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	user, err := getUser(c)
	if err != nil {
		return err
	}

	onConnect := func(conn *gws.Conn) error {
		return Subscribe(c.Request().Context(), conn, scb, user.ID)
	}

	return upgrade(c, onConnect)
}

func (h *handler) GetPages(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
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
	scb, err := getSCB(c)
	if err != nil {
		return err
	}

	if err := c.Request().ParseMultipartForm(2 << 20); err != nil {
		return apiErrors.From(apiErrors.CodeAttachmentSizeExceeded).Wrap(
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

	attachID, err := scb.attachments.Upload(data)
	if err != nil {
		return apiErrors.ErrInternalServerError.Wrap(apiErrors.WithError(err))
	}

	return c.JSON(http.StatusCreated, types.AttachmentResponse{AttachID: attachID})
}

func (h *handler) GetAttachment(c echo.Context) error {
	scb, err := getSCB(c)
	if err != nil {
		return err
	}
	attachID := c.Param("attachId")
	data, MIMEType, err := scb.attachments.Get(attachID)
	if err != nil {
		return apiErrors.ErrNotFound.Wrap(apiErrors.WithError(err))
	}

	return c.Stream(http.StatusOK, MIMEType, data)
}

func getSCB(c echo.Context) (*controlBlock, error) {
	scb, ok := c.Get(SessionCtxKey).(*controlBlock)
	if !ok {
		return nil, echo.ErrForbidden
	}
	return scb, nil
}

func getUser(c echo.Context) (*types.User, error) {
	scb, ok := c.Get(UserCtxKey).(*types.User)
	if !ok {
		return nil, echo.ErrForbidden
	}
	return scb, nil
}
