package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/types"
	"github.com/heat1q/boardsite/session"
	"github.com/heat1q/boardsite/websocket"
)

// Set the api routes
func Set(router *mux.Router) {
	router.HandleFunc("/b/create", handleCreateSession).Methods("POST")
	router.HandleFunc("/b/{id}/users", handleUserCreate).Methods("POST")
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
	idstr := session.Create()

	data := types.CreateBoardResponse{SessionID: idstr}
	json.NewEncoder(w).Encode(data)
}

// handleUserCreate
func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]

	if !session.IsValid(sessionID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userReq := types.User{}
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// new user struct with alias and color
	user, err := session.NewUser(sessionID, userReq.Alias, userReq.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(user)
}

// handleSocketRequest handles request for a websocket upgrade
// based on the sessionID and the userID.
//
// Supported methods: GET
func handleSocketRequest(w http.ResponseWriter, r *http.Request) {
	sessionID, userID := mux.Vars(r)["id"], mux.Vars(r)["userId"]

	if !session.IsValid(sessionID) || !session.IsReadyUser(sessionID, userID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := websocket.UpgradeProtocol(w, r, sessionID, userID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// handlePageRequest handles requests regarding adding or retrieving pages.
//
// Supported methods: GET, POST
func handlePageRequest(w http.ResponseWriter, r *http.Request) {
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
	}
}

// handlePageUpdate handles requests for modifying certain pages.
//
// Supported methods: PUT, DELETE
func handlePageUpdate(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]
	pageID := mux.Vars(r)["pageId"]

	if !session.IsValid(sessionID) || !session.IsValidPage(sessionID, pageID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Fprint(w, session.GetStrokes(sessionID, pageID))
	} else if r.Method == http.MethodPut {
		session.ClearPage(sessionID, pageID)
	} else if r.Method == http.MethodDelete {
		session.DeletePage(sessionID, pageID)
	}
}
