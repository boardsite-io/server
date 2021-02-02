package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/websocket"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// HandleCreateSession handles the request for creating a new session.
// Responds with the unique sessionID of the new session.
//
// Supported methods: POST
func HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	// create new session and set it active
	idstr := session.Create()

	data := types.CreateBoardResponse{ID: idstr}
	json.NewEncoder(w).Encode(data)
}

// HandleSessionRequest handles request for a session based on the sessionID.
//
// Supported methods: GET
func HandleSessionRequest(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]

	if !session.IsValid(sessionID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		if err := websocket.UpgradeProtocol(w, r, sessionID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// HandlePageRequest handles requests regarding adding or retrieving pages.
//
// Supported methods: GET, POST
func HandlePageRequest(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]

	if !session.IsValid(sessionID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		// return pagerank array
		data := types.PageRankResponse{
			PageRank: session.GetPages(sessionID),
		}
		json.NewEncoder(w).Encode(data)
	} else if r.Method == http.MethodPost {
		// add a Page
		data := types.PageRequestData{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} // TODO serialize page data
		session.AddPage(sessionID, data.PageID, data.Index)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// HandlePageUpdate handles requests for modifying certain pages.
//
// Supported methods: PUT, DELETE
func HandlePageUpdate(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]
	pageID := mux.Vars(r)["pageId"]

	if !session.IsValid(sessionID) || !session.IsValidPage(sessionID, pageID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPut {
		session.ClearPage(sessionID, pageID)
	} else if r.Method == http.MethodDelete {
		session.DeletePage(sessionID, pageID)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
