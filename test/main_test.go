package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/heat1q/boardsite/pkg/api"
	"github.com/heat1q/boardsite/pkg/app"
	"github.com/heat1q/boardsite/pkg/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/websocket"
)

const (
	port = 8000
)

var envRedisHost = os.Getenv("BOARDSITE_REDIS_HOST_ENV")

var baseURL = fmt.Sprintf("http://localhost:%d", port)

func TestRunServer(t *testing.T) {
	ctx := context.Background()
	run, shutdown := app.Serve(ctx, port)
	defer shutdown()
	go run()

	client := http.Client{
		Timeout: time.Second,
	}

	// create board
	r, err := client.Post(baseURL+"/board/create", "application/json", strings.NewReader(""))
	assert.NoError(t, err)

	boardResp := api.CreateBoardResponse{}
	json.NewDecoder(r.Body).Decode(&boardResp)
	assert.NotEqual(t, boardResp.ID, "", "id not in create response")
	sessionID := boardResp.ID

	// start websocket
	sessionURL := fmt.Sprintf("ws://localhost:%d/board/%s", port, sessionID)
	conn, wsResp, wsErr := websocket.DefaultDialer.Dial(sessionURL, nil)
	require.NoError(t, wsErr)
	require.Equal(t, http.StatusSwitchingProtocols, wsResp.StatusCode)
	defer conn.Close()

	// initial data on fresh session should be empty
	msgType, msg, msgErr := conn.ReadMessage()
	require.NoError(t, msgErr)
	require.Equal(t, msgType, websocket.TextMessage)
	require.Equal(t, msg, []byte("[]"))

	// send some data
	testStroke1 := `{"id":"testid1","type":"line","color":"#ff00ff","line_width":1,"position":[1.3,3.7]}`
	testStroke2 := `{"id":"testid2","type":"rect","color":"#0000ff","line_width":2,"position":[4.2,1.5]}`
	strokeData := fmt.Sprintf("[%s,%s]", testStroke1, testStroke2)
	wsWErr := conn.WriteMessage(websocket.TextMessage, []byte(strokeData))
	require.NoError(t, wsWErr)

	// delete stroke 2
	conn.WriteMessage(websocket.TextMessage, []byte(`[{"id":"testid2","type":"delete"}]`))

	// database connection
	db, dbErr := database.NewRedisConn(sessionID)
	defer db.Close()
	require.NoError(t, dbErr)

	// check the database entries
	time.Sleep(time.Second / 2)
	strokeDb, fetchErr := db.FetchAll()
	require.NoError(t, fetchErr)
	require.Equal(t, fmt.Sprintf("[%s]", testStroke1), strokeDb)

	// clear the board
	reqClr, _ := http.NewRequest("PUT", baseURL+"/board/"+sessionID, strings.NewReader(`{"action":"clear"}`))
	respClr, errClr := client.Do(reqClr)
	require.NoError(t, errClr)
	require.Equal(t, http.StatusOK, respClr.StatusCode)

	// check if db is clear
	time.Sleep(time.Second / 2)
	strokeDb2, fetchErr2 := db.FetchAll()
	require.NoError(t, fetchErr2)
	require.Equal(t, "[]", strokeDb2)
}
