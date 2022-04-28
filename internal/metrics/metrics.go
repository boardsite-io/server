package metrics

import (
	"net/http"
	"strconv"

	echoProm "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/boardsite-io/server/internal/session"
)

type Handler interface {
	GetMetrics(c echo.Context) error
	MiddlewareFunc() echo.MiddlewareFunc
}

type handler struct {
	prom        *echoProm.Prometheus
	promHandler http.Handler
	dispatcher  session.Dispatcher
	numReq      *echoProm.Metric
	numSessions *echoProm.Metric
	numUsers    *echoProm.Metric
}

func NewHandler(dispatcher session.Dispatcher) Handler {
	promhttp.Handler()
	skipper := func(echo.Context) bool { return true }
	h := handler{
		dispatcher: dispatcher,
		numReq: &echoProm.Metric{
			ID:          "numReq",
			Name:        "num_requests_total",
			Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
			Type:        "counter_vec",
			Args:        []string{"code", "method", "host", "url"}},
		numSessions: &echoProm.Metric{
			ID:          "numSessions",
			Name:        "num_sessions",
			Description: "The number of active session.",
			Type:        "gauge"},
		numUsers: &echoProm.Metric{
			ID:          "numUsers",
			Name:        "num_users",
			Description: "The total number of active users across all sessions.",
			Type:        "gauge"},
	}
	h.prom = echoProm.NewPrometheus(
		"boardsite",
		skipper,
		[]*echoProm.Metric{h.numReq, h.numSessions, h.numUsers})
	h.promHandler = promhttp.Handler()
	return &h
}

func (h *handler) GetMetrics(c echo.Context) error {
	numSession, ok := h.numSessions.MetricCollector.(prometheus.Gauge)
	if ok {
		numSession.Set(float64(h.dispatcher.NumSessions()))
	}

	numUsers, ok := h.numUsers.MetricCollector.(prometheus.Gauge)
	if ok {
		numUsers.Set(float64(h.dispatcher.NumUsers()))
	}

	h.promHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *handler) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			numReq, ok := h.numReq.MetricCollector.(*prometheus.CounterVec)
			if ok {
				numReq.WithLabelValues(
					strconv.Itoa(c.Response().Status),
					c.Request().Method,
					c.RealIP(),
					c.Path()).Inc()
			}

			return err
		}
	}
}
