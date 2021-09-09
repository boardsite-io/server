package routes

import (
	"encoding/json"
	"net/http"

	"github.com/heat1q/boardsite/api/types"
	apiErrors "github.com/heat1q/boardsite/api/types/errors"
)

type requestContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *requestContext) Request() *http.Request {
	return c.r
}

func (c *requestContext) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *requestContext) JSON(status int, v interface{}) error {
	// wrap the content in message
	c.w.Header().Add("Content-Type", "application/json")
	c.w.WriteHeader(status)
	msg := types.NewMessage(v, "")
	return json.NewEncoder(c.w).Encode(msg)
}

func (c *requestContext) Stream(status int, data []byte, MIMEType string) error {
	c.w.Header().Add("Content-Type", MIMEType)
	c.w.WriteHeader(status)
	_, err := c.w.Write(data)
	return err
}

func (c *requestContext) NoContent(status int) error {
	c.w.WriteHeader(status)
	return nil
}

func handleRequest(fn func(c *requestContext) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &requestContext{w: w, r: r}
		err := fn(c)
		if err != nil {
			err = apiErrors.MaptoHTTPError(err)
			c.w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(err.(*apiErrors.Wrapper).Code)
			_ = json.NewEncoder(w).Encode(err)
		}
	}
}
