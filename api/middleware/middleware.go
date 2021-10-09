package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/heat1q/boardsite/api/request"
	apiErrors "github.com/heat1q/boardsite/api/types/errors"
)

func RequestLogger(next request.HandlerFunc) request.HandlerFunc {
	return func(c *request.Context) error {
		// measure response time
		start := time.Now()
		reqBody, _ := io.ReadAll(c.RequestBody())

		err := next(c)

		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("-- incoming request (%d) %s %s -- %s",
			c.StatusCode(),
			c.Request().Method,
			c.Request().RequestURI,
			c.Request().Header.Get("User-Agent"),
		))

		if len(reqBody) > 0 {
			var body bytes.Buffer
			if err := json.Compact(&body, reqBody); err == nil {
				sb.WriteString(fmt.Sprintf(" -- req body: %s", body.String()))
			} else {
				sb.WriteString(fmt.Sprintf(" -- req body content length: %d", len(reqBody)))
			}
		}

		respBody, _ := io.ReadAll(c.ResponseBody())
		if len(respBody) > 0 {
			// dont spam the logs with huge responses
			if len(respBody) < 2<<10 && json.Valid(respBody) {
				sb.WriteString(fmt.Sprintf(" -- resp body: %s", string(respBody)))
			} else {
				sb.WriteString(fmt.Sprintf(" -- resp body content length: %d", len(respBody)))
			}
		}

		if err != nil {
			sb.WriteString(fmt.Sprintf(" -- error: %v", err))
		}

		elapsed := time.Since(start)
		sb.WriteString(fmt.Sprintf(" -- took %s", elapsed))

		log.Println(sb.String())

		return err
	}
}

func ErrorMapper(next request.HandlerFunc) request.HandlerFunc {
	return func(c *request.Context) error {
		err := next(c)

		// write error response
		if err != nil {
			err = apiErrors.MaptoHTTPError(err)
			_ = c.JSON(err.(*apiErrors.Wrapper).StatusCode, err)
		}

		return err
	}
}
