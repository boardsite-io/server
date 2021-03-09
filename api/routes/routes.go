package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/websocket"
)

// Set the api routes
func Set(router *mux.Router) {
	router.HandleFunc("/b/create", handleCreateSession).Methods("POST")
	router.HandleFunc("/b/{id}/users", handleUsers).Methods("GET", "POST")
	router.HandleFunc("/b/{id}/users/{userId}/socket", handleSocketRequest).Methods("GET")
	router.HandleFunc("/b/{id}/pages", handlePageRequest).Methods("GET", "POST")
	router.HandleFunc("/b/{id}/pages/{pageId}", handlePageUpdate).Methods("GET", "PUT", "DELETE")
}

// handleCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
//
// Supported methods: POST
func handleCreateSession(w http.ResponseWriter, r *http.Request) {
	// create new session and set it active
	idstr, err := session.Create()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeMessage(w, types.NewMessage(idstr, ""))
}

// handleUserCreate
func handleUsers(w http.ResponseWriter, r *http.Request) {
	scb, err := session.GetSCB(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	if r.Method == http.MethodGet {
		writeMessage(w, types.NewMessage(scb.GetUsers(), ""))
	} else if r.Method == http.MethodPost {
		var userReq types.User
		if err := types.DecodeMsgContent(r.Body, &userReq); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		// new user struct with alias and color
		user, err := session.NewUser(scb, userReq.Alias, userReq.Color)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeMessage(w, types.NewMessage(user, ""))
	}
}

// handleSocketRequest handles request for a websocket upgrade
// based on the sessionID and the userID.
//
// Supported methods: GET
func handleSocketRequest(w http.ResponseWriter, r *http.Request) {
	sessionID, userID := mux.Vars(r)["id"], mux.Vars(r)["userId"]

	scb, err := session.GetSCB(sessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	_, errUser := scb.GetUserReady(userID)
	if errUser != nil {
		writeError(w, http.StatusNotFound, errUser)
		return
	}

	if err := websocket.UpgradeProtocol(w, r, scb, userID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}

// handlePageRequest handles requests regarding adding or retrieving pages.
//
// Supported methods: GET, POST
func handlePageRequest(w http.ResponseWriter, r *http.Request) {
	scb, err := session.GetSCB(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	if r.Method == http.MethodGet {
		pageRank, meta, err := session.GetPages(scb.ID)
		if err != nil {
			writeError(w, http.StatusServiceUnavailable, err)
		}

		// return pagerank array
		writeMessage(w, types.NewMessage(types.ContentPageSync{
			PageRank: pageRank,
			Meta:     meta,
		}, ""))
	} else if r.Method == http.MethodPost {
		// add a Page
		var data types.ContentPageRequest
		if err := types.DecodeMsgContent(r.Body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := session.AddPage(scb, data.PageID, data.Index, &data.PageMeta); err != nil {
			writeError(w, http.StatusServiceUnavailable, err)
		}
	}
}

// handlePageUpdate handles requests for modifying certain pages.
//
// Supported methods: PUT, DELETE
func handlePageUpdate(w http.ResponseWriter, r *http.Request) {
	scb, err := session.GetSCB(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	pageID := mux.Vars(r)["pageId"]
	if !session.IsValidPage(scb.ID, pageID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		strokes, errFetch := session.GetStrokes(scb.ID, pageID)
		if errFetch != nil {
			writeError(w, http.StatusServiceUnavailable, errFetch)
		}
		writeMessage(
			w,
			types.NewMessage(strokes, ""),
		)
	} else if r.Method == http.MethodPut {
		if err := session.ClearPage(scb, pageID); err != nil {
			writeError(w, http.StatusServiceUnavailable, err)
		}
	} else if r.Method == http.MethodDelete {
		if err := session.DeletePage(scb, pageID); err != nil {
			writeError(w, http.StatusServiceUnavailable, err)
		}
	}
}

func writeMessage(w http.ResponseWriter, content interface{}) {
	if err := json.NewEncoder(w).Encode(content); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(types.NewErrorMessage(err))
}
