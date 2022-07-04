package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/boardsite-io/server/pkg/constant"
	"github.com/boardsite-io/server/pkg/log"
)

func CORS(allowedOrigins string) echo.MiddlewareFunc {
	origins := strings.Split(allowedOrigins, ",")
	log.Global().Infof("CORS: allowed origins: %v", origins)
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization, constant.HeaderUserID, constant.HeaderSessionSecret},
	})
}
