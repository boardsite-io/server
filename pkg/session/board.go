package session

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/pkg/api"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// CreateBoard creates a new board with parameters X and Y and redirects
// to "/board/{id}" by setting a unique ID.
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// TODO retrieve x,y from form
	form := api.SetupForm{}
	// if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
	// 	return
	// }

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

	// assign to SessionControl struct
	ActiveSession[idstr] = NewSessionControl(idstr, form.X, form.Y)

	data := api.CreateBoardResponse{ID: idstr}
	json.NewEncoder(w).Encode(data)
}

// HandleBoardRequest starts the websocket based on route "/board/{id}"
// if a session with {id} has been create, i.e. is active.
func HandleBoardRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// session does not exist
	if ActiveSession[vars["id"]] == nil {
		// TODO return status 404
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == "PUT" {
		// modify session
		data := api.BoardRequest{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Action == "clear" {
			ActiveSession[vars["id"]].Clear()
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	} else if r.Method == "GET" {
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
}
