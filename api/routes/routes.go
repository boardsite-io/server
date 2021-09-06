package routes

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/types"
	apiErrors "github.com/heat1q/boardsite/api/types/errors"
	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/websocket"
)

// Set the api routes
func Set(router *mux.Router) {
	router.HandleFunc("/b/create", handleRequest(postCreateSession)).Methods(http.MethodPost)
	router.HandleFunc("/b/{id}/users", handleRequest(getUsers)).Methods(http.MethodGet)
	router.HandleFunc("/b/{id}/users", handleRequest(postUsers)).Methods(http.MethodPost)
	router.HandleFunc("/b/{id}/users/{userId}/socket", handleRequest(getSocket)).Methods(http.MethodGet)
	router.HandleFunc("/b/{id}/pages", handleRequest(getPages)).Methods(http.MethodGet)
	router.HandleFunc("/b/{id}/pages", handleRequest(postPages)).Methods(http.MethodPost)
	router.HandleFunc("/b/{id}/pages", handleRequest(putPages)).Methods(http.MethodPut)
	router.HandleFunc("/b/{id}/pages", handleRequest(deletePages)).Methods(http.MethodDelete)
	router.HandleFunc("/b/{id}/pages/{pageId}", handleRequest(getPageUpdate)).Methods(http.MethodGet)
	router.HandleFunc("/b/{id}/pages/{pageId}", handleRequest(deletePageUpdate)).Methods(http.MethodDelete)
	router.HandleFunc("/b/{id}/attachments", handleRequest(postAttachment)).Methods(http.MethodPost)
	router.HandleFunc("/b/{id}/attachments/{attachId}", handleRequest(getAttachment)).Methods(http.MethodGet)
}

// postCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
func postCreateSession(c *requestContext) error {
	idstr, err := session.Create()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, idstr)
}

func getUsers(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound
	}
	return c.JSON(http.StatusOK, scb.GetUsers())
}

// handleUserCreate
func postUsers(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound
	}

	var userReq types.User
	if err := types.DecodeMsgContent(c.Request().Body, &userReq); err != nil {
		return apiErrors.BadRequest
	}

	// new user struct with alias and color
	user, err := session.NewUser(scb, userReq.Alias, userReq.Color)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// getSocket handles request for a websocket upgrade
// based on the sessionID and the userID.
func getSocket(c *requestContext) error {
	sessionID, userID := mux.Vars(c.Request())["id"], mux.Vars(c.Request())["userId"]

	scb, err := session.GetSCB(sessionID)
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}
	_, errUser := scb.GetUserReady(userID)
	if errUser != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	return websocket.UpgradeProtocol(c.ResponseWriter(), c.Request(), scb, userID)
}

func getPages(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageRank, meta, err := session.GetPages(scb.ID)
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
func postPages(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	// add a Page
	var data types.ContentPageRequest
	if err := types.DecodeMsgContent(c.Request().Body, &data); err != nil {
		return apiErrors.BadRequest.SetInfo(err)
	}

	if err := session.AddPages(scb, data.PageID, data.Index, data.Meta); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusCreated)
}

func putPages(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	var data types.ContentPageRequest
	if err := types.DecodeMsgContent(c.Request().Body, &data); err != nil {
		return apiErrors.BadRequest
	}

	if err := session.UpdatePages(scb, data.PageID, data.Meta, data.Clear); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func deletePages(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	var data types.ContentPageRequest
	if err := types.DecodeMsgContent(c.Request().Body, &data); err != nil {
		return apiErrors.BadRequest
	}

	if err := session.DeletePages(scb, data.PageID...); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func getPageUpdate(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageID := mux.Vars(c.Request())["pageId"]
	if !session.IsValidPage(scb.ID, pageID) {
		return apiErrors.NotFound
	}

	strokes, errFetch := session.GetStrokes(scb.ID, pageID)
	if errFetch != nil {
		return apiErrors.InternalServerError.SetInfo(errFetch)
	}

	return c.JSON(http.StatusOK, strokes)
}

func deletePageUpdate(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	pageID := mux.Vars(c.Request())["pageId"]
	if !session.IsValidPage(scb.ID, pageID) {
		return apiErrors.NotFound
	}

	if err := session.DeletePages(scb, pageID); err != nil {
		return apiErrors.InternalServerError.SetInfo(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func postAttachment(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
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

	return c.JSON(http.StatusCreated, attachID)
}

func getAttachment(c *requestContext) error {
	scb, err := session.GetSCB(mux.Vars(c.Request())["id"])
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}
	attachID := mux.Vars(c.Request())["attachId"]
	data, MIMEType, err := scb.Attachments.Get(attachID)
	if err != nil {
		return apiErrors.NotFound.SetInfo(err)
	}

	return c.Stream(http.StatusOK, data, MIMEType)
}
