package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/boardsite-io/server/pkg/constant"
	"github.com/boardsite-io/server/pkg/log"
)

var loggedHeaders = map[string]struct{}{
	echo.HeaderXForwardedFor: {},
	echo.HeaderContentType:   {},
	echo.HeaderContentLength: {},
	echo.HeaderOrigin:        {},
	"User-Agent":             {},
	constant.HeaderUserID:    {},
}

var jsonBody = regexp.MustCompile("^application/json.*")

func RequestLogger() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// measure response time
			start := time.Now()
			meta := make(map[string]any, 32)

			var reqBody, respBody []byte
			cfg := middleware.BodyDumpConfig{
				Handler: func(c echo.Context, req []byte, resp []byte) {
					reqBody = req
					respBody = resp
				},
				Skipper: func(c echo.Context) bool {
					return strings.Contains(c.Request().RequestURI, "socket")
				},
			}
			bodyDump := middleware.BodyDumpWithConfig(cfg)

			err := bodyDump(next)(c)

			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("incoming request (%d) %s %s",
				c.Response().Status,
				c.Request().Method,
				c.Request().RequestURI,
			))

			meta["Req.HttpMethod"] = c.Request().Method
			meta["Req.Path"] = c.Path()
			setMeta(c, reqBody, meta, "Req")
			setMeta(c, respBody, meta, "Resp")

			elapsed := time.Since(start)
			meta["duration"] = elapsed.String()

			if err != nil {
				meta["error"] = err.Error()
				log.Ctx(c.Request().Context()).With(log.WithMeta(meta)...).Error(sb.String())
			} else {
				log.Ctx(c.Request().Context()).With(log.WithMeta(meta)...).Info(sb.String())
			}

			return err
		}
	}
}

func setMeta(c echo.Context, body []byte, meta map[string]any, prefix string) {
	for k, v := range c.Response().Header() {
		if _, ok := loggedHeaders[k]; ok {
			meta[fmt.Sprintf("%s.Header.%s", prefix, k)] = strings.Join(v, ";")
		}
	}
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if len(body) > 0 && len(body) < 2<<10 && jsonBody.MatchString(contentType) {
		var buf bytes.Buffer
		_ = json.Compact(&buf, body)
		meta[prefix+".Body"] = buf.String()
	}
}
