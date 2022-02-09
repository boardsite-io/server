package api

import (
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/heat1q/boardsite/api/middleware"
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b", echomw.Gzip(), middleware.RequestLogger())
	boardGroup.POST( /*  */ "/create", s.session.PostCreateSession, middleware.RateLimiting(s.cfg.Session.RPM, middleware.WithIP()))

	usersGroup := boardGroup.Group("/:id/users")
	usersGroup.POST( /*  */ "", s.session.PostUsers)
	usersGroup.GET( /*   */ "", s.session.GetUsers, middleware.Session(s.dispatcher))
	usersGroup.GET( /*   */ "/:userId/socket", s.session.GetSocket, middleware.Session(s.dispatcher))

	pagesGroup := boardGroup.Group("/:id/pages", middleware.Session(s.dispatcher))
	pagesGroup.GET( /*   */ "", s.session.GetPageRank)
	pagesGroup.POST( /*  */ "", s.session.PostPages)
	pagesGroup.PUT( /*   */ "", s.session.PutPages)
	pagesGroup.GET( /*   */ "/:pageId", s.session.GetPage)
	pagesGroup.GET( /*   */ "/sync", s.session.GetPageSync)
	pagesGroup.POST( /*  */ "/sync", s.session.PostPageSync)

	attachGroup := boardGroup.Group("/:id/attachments", middleware.Session(s.dispatcher))
	attachGroup.POST( /**/ "", s.session.PostAttachment, middleware.RateLimiting(s.cfg.Session.RPM, middleware.WithUserIP()))
	attachGroup.GET( /* */ "/:attachId", s.session.GetAttachment)

	if s.cfg.Server.Metrics.Enabled {
		s.setMetricsRoutes()
	}
}

func (s *Server) setMetricsRoutes() {
	metricsGroup := s.echo.Group(s.cfg.Server.Metrics.Route,
		middleware.BasicAuth(s.cfg.Server.Metrics.User, s.cfg.Server.Metrics.Password))
	metricsGroup.GET("", s.metrics.GetMetrics)
}
