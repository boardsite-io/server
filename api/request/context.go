package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request

	statusCode int
	Headers    http.Header
	body       []byte
}

type HandlerFunc func(*Context) error

// NewContext creates a new request context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:       w,
		r:       r,
		Headers: http.Header{},
	}
}

// Ctx returns the underlying context of the current request.
func (c *Context) Ctx() context.Context {
	return c.r.Context()
}

func (c *Context) Request() *http.Request {
	return c.r
}

// RequestBody returns a safe to use io.Reader to read the request body
func (c *Context) RequestBody() io.Reader {
	bodyDump, _ := io.ReadAll(c.Request().Body)
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyDump)) // reset the request body
	return bytes.NewBuffer(bodyDump)
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *Context) StatusCode() int {
	return c.statusCode
}

func (c *Context) ResponseBody() io.Reader {
	return bytes.NewBuffer(c.body)
}

// Vars returns the route variables for the current request, if any.
func (c *Context) Vars() map[string]string {
	return mux.Vars(c.r)
}

func (c *Context) JSON(status int, v interface{}) error {
	// wrap the content in message
	c.Headers.Add("Content-Type", "application/json")
	c.statusCode = status
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.body = b
	return nil
}

func (c *Context) Stream(status int, data []byte, MIMEType string) error {
	c.Headers.Add("Content-Type", MIMEType)
	c.statusCode = status
	c.body = data
	return nil
}

func (c *Context) NoContent(status int) error {
	c.statusCode = status
	return nil
}

// NewHandler wraps a Context handler and returns an http handler functions
func NewHandler(handl HandlerFunc, mwFn ...func(fn HandlerFunc) HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		// apply middlewares
		next := handl
		for _, fn := range mwFn {
			next = fn(next)
		}

		// error should be resolved at this point
		if err := next(c); err != nil && !isErrorStatusCode(c.statusCode) {
			log.Printf("unhandeld error: %v", err)
		}

		// write response
		for k, vals := range c.Headers {
			for _, v := range vals {
				c.w.Header().Add(k, v)
			}
		}
		c.w.WriteHeader(c.statusCode)
		if _, err := c.w.Write(c.body); err != nil {
			c.w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func isErrorStatusCode(status int) bool {
	return status/100 == 4 || status/100 == 5
}
