package session

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"boardsite/api/board"
	"boardsite/api/database"

	"github.com/gorilla/mux"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	// ActiveSession maps the session is to the SessionControl struct
	ActiveSession = make(map[string]*board.SessionControl)
)

type createResponse struct {
	ID string `json:"id"`
}

// CreateBoard creates a new board with parameters X and Y and redirects
// to "/board/{id}" by setting a unique ID.
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	form := board.SetupForm{}
	// TODO retrieve x,y from form
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 6)
	// find available id
	for {
		for i := range id {
			id[i] = letters[rand.Intn(len(letters))]
		}

		if ActiveSession[string(id)] == nil {
			break
		}
	}
	idstr := string(id)

	db, err := database.NewConnection(idstr)
	if err != nil {
		return
	}

	// assign to SessionControl struct
	ActiveSession[idstr] = board.NewSessionControl(idstr, form.X, form.Y, db)

	data := createResponse{ID: idstr}
	json.NewEncoder(w).Encode(data)
}

// ServeBoard starts the websocket based on route "/board/{id}"
// if a session with {id} has been create, i.e. is active.
func ServeBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// session does not exist
	if ActiveSession[vars["id"]] == nil || !ActiveSession[vars["id"]].IsActive {
		// TODO return status 404
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// upgrade to websocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	onClientConnect(vars["id"], conn)
	defer onClientDisconnect(vars["id"], conn)

	InitWebsocket(vars["id"], conn)
}

func closeSession(sessionID string) {
	ActiveSession[sessionID].IsActive = false
	ActiveSession[sessionID].DB.Clear()
	delete(ActiveSession, sessionID)
}
