package middleware

import (
	"net/http"
	"strings"

	"github.com/heat1q/boardsite/api/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS(allowedOrigins string) echo.MiddlewareFunc {
	origins := strings.Split(allowedOrigins, ",")
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization, types.HeaderUserID, types.HeaderSessionSecret},
	})
}
