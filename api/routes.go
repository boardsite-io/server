package api

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/heat1q/boardsite/api/middleware"
)

const (
	rpmSession = 100
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	metricsGroup := s.echo.Group(s.cfg.Server.Metrics.Route,
		middleware.BasicAuth(s.cfg.Server.Metrics.User, s.cfg.Server.Metrics.Password))
	metricsGroup.GET("", prometheusHandler())

	boardGroup := s.echo.Group("/b", echomw.Gzip(), middleware.RequestLogger())
	boardGroup.POST( /*  */ "/create", s.session.PostCreateSession, middleware.RateLimiting(1, middleware.WithIP()))

	usersGroup := boardGroup.Group("/:id/users")
	usersGroup.POST( /*  */ "", s.session.PostUsers)
	usersGroup.GET( /*   */ "", s.session.GetUsers, middleware.Session(s.dispatcher))
	usersGroup.GET( /*   */ "/:userId/socket", s.session.GetSocket, middleware.Session(s.dispatcher))

	pagesGroup := boardGroup.Group("/:id/pages", middleware.Session(s.dispatcher))
	pagesGroup.GET( /*   */ "", s.session.GetPages) // get page rank
	pagesGroup.POST( /*  */ "", s.session.PostPages)
	pagesGroup.PUT( /*   */ "", s.session.PutPages)
	pagesGroup.DELETE( /**/ "", s.session.DeletePages)
	pagesGroup.GET( /*   */ "/:pageId", s.session.GetPageUpdate)
	pagesGroup.DELETE( /**/ "/:pageId", s.session.DeletePageUpdate)

	attachGroup := boardGroup.Group("/:id/attachments", middleware.Session(s.dispatcher))
	attachGroup.POST( /**/ "", s.session.PostAttachment, middleware.RateLimiting(1, middleware.WithUserIP()))
	attachGroup.GET( /* */ "/:attachId", s.session.GetAttachment)
}

func prometheusHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
