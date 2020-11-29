package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/heat1q/boardsite/pkg/api"
	"github.com/heat1q/boardsite/pkg/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/websocket"
)

const (
	port = 8000
)

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

	// wsWErr := conn.WriteMessage(websocket.TextMessage, []byte("[]"))
	// assert.NoError(t, wsWErr)
}
