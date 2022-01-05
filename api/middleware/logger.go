package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/heat1q/boardsite/api/log"
)

func RequestLogger() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// measure response time
			start := time.Now()
			meta := make(map[string]interface{}, 32)

			var reqBody, respBody []byte
			cfg := middleware.BodyDumpConfig{
				Handler: func(c echo.Context, req []byte, resp []byte) {
					reqBody = req
					respBody = resp
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

			setRequestMeta(c, reqBody, meta)
			setResponseMeta(c, respBody, meta)

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

func setRequestMeta(c echo.Context, reqBody []byte, meta map[string]interface{}) {
	meta["Req.HttpMethod"] = c.Request().Method
	meta["Req.Path"] = c.Request().RequestURI

	for k, v := range c.Request().Header {
		meta[fmt.Sprintf("Req.Header.%s", k)] = strings.Join(v, ";")
	}

	if len(reqBody) > 0 {
		var body bytes.Buffer
		if err := json.Compact(&body, reqBody); err == nil {
			meta["Req.Body"] = body.String()
		} else {
			meta["Req.ContentLength"] = c.Request().ContentLength
		}
	}
}

func setResponseMeta(c echo.Context, respBody []byte, meta map[string]interface{}) {
	for k, v := range c.Response().Header() {
		meta[fmt.Sprintf("Resp.Header.%s", k)] = strings.Join(v, ";")
	}

	if len(respBody) > 0 {
		// dont spam the logs with huge responses
		if len(respBody) < 2<<10 && json.Valid(respBody) {
			var body bytes.Buffer
			_ = json.Compact(&body, respBody)
			meta["Resp.Body"] = body.String()
		} else {
			meta["Resp.ContentLength"] = c.Request().Response.ContentLength
		}
	}
}
