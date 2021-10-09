package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/middleware"
	"github.com/heat1q/boardsite/api/request"
	"github.com/heat1q/boardsite/api/types"
	apiErrors "github.com/heat1q/boardsite/api/types/errors"
	"github.com/heat1q/boardsite/websocket"
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	s.setHandleFunc("/b/create", s.postCreateSession).Methods(http.MethodPost)
	s.setHandleFunc("/b/{id}/users", s.getUsers).Methods(http.MethodGet)
	s.setHandleFunc("/b/{id}/users", s.postUsers).Methods(http.MethodPost)
	s.setHandleFunc("/b/{id}/users/{userId}/socket", s.getSocket).Methods(http.MethodGet)
	s.setHandleFunc("/b/{id}/pages", s.getPages).Methods(http.MethodGet)
	s.setHandleFunc("/b/{id}/pages", s.postPages).Methods(http.MethodPost)
	s.setHandleFunc("/b/{id}/pages", s.putPages).Methods(http.MethodPut)
	s.setHandleFunc("/b/{id}/pages", s.deletePages).Methods(http.MethodDelete)
	s.setHandleFunc("/b/{id}/pages/{pageId}", s.getPageUpdate).Methods(http.MethodGet)
	s.setHandleFunc("/b/{id}/pages/{pageId}", s.deletePageUpdate).Methods(http.MethodDelete)
	s.setHandleFunc("/b/{id}/attachments", s.postAttachment).Methods(http.MethodPost)
	s.setHandleFunc("/b/{id}/attachments/{attachId}", s.getAttachment).Methods(http.MethodGet)
}

func (s *Server) setHandleFunc(path string, fn request.HandlerFunc) *mux.Route {
	return s.router.HandleFunc(path, request.NewHandler(fn, middleware.ErrorMapper, middleware.RequestLogger))
}

// postCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
func (s *Server) postCreateSession(c *request.Context) error {
	idstr, err := s.dispatcher.Create(s.cfg.Session.MaxUsers)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, types.CreateSessionResponse{SessionID: idstr})
}

func (s *Server) getUsers(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}
	return c.JSON(http.StatusOK, scb.GetUsers())
}

// handleUserCreate
func (s *Server) postUsers(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	var userReq types.User

	if err := json.NewDecoder(c.Request().Body).Decode(&userReq); err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	// new user struct with alias and color
	user, err := scb.NewUser(userReq.Alias, userReq.Color)
	if err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	return c.JSON(http.StatusCreated, user)
}

// getSocket handles request for a websocket upgrade
// based on the sessionID and the userID.
func (s *Server) getSocket(c *request.Context) error {
	var (
		sessionID = c.Vars()["id"]
		userID    = c.Vars()["userId"]
	)

	scb, err := s.dispatcher.GetSCB(sessionID)
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}
	_, errUser := scb.GetUserReady(userID)
	if errUser != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	return websocket.UpgradeProtocol(c.Ctx(), c.ResponseWriter(), c.Request(), scb, userID)
}

func (s *Server) getPages(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageRank, meta, err := scb.GetPages(c.Ctx())
	if err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	// return pagerank array
	pages := types.ContentPageSync{
		PageRank: pageRank,
		Meta:     meta,
	}
	return c.JSON(http.StatusOK, pages)
}

// handlePageRequest handles requests regarding adding or retrieving pages.
func (s *Server) postPages(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	// add a Page
	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	if err := scb.AddPages(c.Ctx(), data.PageID, data.Index, data.Meta); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) putPages(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	if err := scb.UpdatePages(c.Ctx(), data.PageID, data.Meta, data.Clear); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deletePages(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	var data types.ContentPageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	if err := scb.DeletePages(c.Ctx(), data.PageID...); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) getPageUpdate(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageID := c.Vars()["pageId"]
	if !scb.IsValidPage(c.Ctx(), pageID) {
		return apiErrors.NotFound
	}

	strokes, errFetch := scb.GetStrokes(c.Ctx(), pageID)
	if errFetch != nil {
		return apiErrors.InternalServerError.SetInfo(errFetch)
	}

	return c.JSON(http.StatusOK, strokes)
}

func (s *Server) deletePageUpdate(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageID := c.Vars()["pageId"]
	if !scb.IsValidPage(c.Ctx(), pageID) {
		return apiErrors.NotFound
	}

	if err := scb.DeletePages(c.Ctx(), pageID); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) postAttachment(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	if err := c.Request().ParseMultipartForm(2 << 20); err != nil {
		return apiErrors.BadRequest.SetInfo("file size exceeded limit of 2MB")
	}
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	attachID, err := scb.Attachments.Upload(data)
	if err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.JSON(http.StatusCreated, types.AttachmentResponse{AttachID: attachID})
}

func (s *Server) getAttachment(c *request.Context) error {
	scb, err := s.dispatcher.GetSCB(c.Vars()["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}
	attachID := c.Vars()["attachId"]
	data, MIMEType, err := scb.Attachments.Get(attachID)
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	return c.Stream(http.StatusOK, data, MIMEType)
}
