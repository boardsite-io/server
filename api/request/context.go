package request

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/heat1q/boardsite/api/types"
	apiErrors "github.com/heat1q/boardsite/api/types/errors"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

// NewContext creates a new request context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
	}
}

// Ctx returns the underlying context of the current request.
func (c *Context) Ctx() context.Context {
	return c.r.Context()
}

func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

// Vars returns the route variables for the current request, if any.
func (c *Context) Vars() map[string]string {
	return mux.Vars(c.r)
}

func (c *Context) JSON(status int, v interface{}) error {
	// wrap the content in message
	c.w.Header().Add("Content-Type", "application/json")
	c.w.WriteHeader(status)
	msg := types.NewMessage(v, "")
	return json.NewEncoder(c.w).Encode(msg)
}

func (c *Context) Stream(status int, data []byte, MIMEType string) error {
	c.w.Header().Add("Content-Type", MIMEType)
	c.w.WriteHeader(status)
	_, err := c.w.Write(data)
	return err
}

func (c *Context) NoContent(status int) error {
	c.w.WriteHeader(status)
	return nil
}

// NewHandler wraps a Context handler and returns an http handler functions
func NewHandler(fn func(c *Context) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		err := fn(c)
		if err != nil {
			err = apiErrors.MaptoHTTPError(err)
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			w.WriteHeader(err.(*apiErrors.Wrapper).Code)
			_ = json.NewEncoder(w).Encode(err)
		}
	}
}
