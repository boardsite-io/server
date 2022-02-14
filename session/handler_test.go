package session_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/heat1q/boardsite/api/config"
	"github.com/heat1q/boardsite/session"
	"github.com/stretchr/testify/assert"

	"github.com/heat1q/boardsite/session/sessionfakes"

	"github.com/labstack/echo/v4"
)

func Test_handler_PostCreateWithConfig(t *testing.T) {
	const sid = "sid"
	e := echo.New()
	dispatcher := &sessionfakes.FakeDispatcher{}
	cfg := &config.Session{
		MaxUsers: 50,
	}
	handler := session.NewHandler(cfg, dispatcher)

	tests := []struct {
		name     string
		reqBody  string
		wantBody string
		wantErr  bool
	}{
		{
			name:     "successful with valid config",
			reqBody:  `{"maxUsers": 20}`,
			wantBody: `{"sessionId": "sid"}`,
		},
		{
			name:    "fails with invalid maxUsers",
			reqBody: `{"maxUsers": -1}`,
			wantErr: true,
		},
		{
			name:    "fails with invalid maxUsers",
			reqBody: `{"maxUsers": 0}`,
			wantErr: true,
		},
		{
			name:    "fails with invalid maxUsers",
			reqBody: `{"maxUsers": 100}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher.CreateReturns(sid, nil)
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			r.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rr := httptest.NewRecorder()
			c := e.NewContext(r, rr)

			err := handler.PostCreateWithConfig(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
