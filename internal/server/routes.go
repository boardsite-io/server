package server

import (
	echomw "github.com/labstack/echo/v4/middleware"

	apimw "github.com/boardsite-io/server/internal/middleware"
	libmw "github.com/boardsite-io/server/pkg/middleware"
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b", echomw.Gzip(), libmw.RequestLogger())

	createGroup := boardGroup.Group("/create", libmw.RateLimiting(s.cfg.Server.RPM, libmw.WithIP()))
	createGroup.POST( /**/ "", s.session.PostCreateSession)
	createGroup.POST( /**/ "/config", s.session.PostCreateSessionConfig)

	configGroup := boardGroup.Group("/:id/config", apimw.Session(s.dispatcher))
	configGroup.GET( /*  */ "", s.session.GetSessionConfig)

	hostGroup := boardGroup.Group("", apimw.Session(s.dispatcher), apimw.Host())
	hostGroup.PUT("/:id/config", s.session.PutSessionConfig)
	hostGroup.PUT("/:id/users/:userId", s.session.PutKickUser)

	usersGroup := boardGroup.Group("/:id/users")
	usersGroup.POST( /* */ "", s.session.PostUsers)
	usersGroup.GET( /*  */ "/:userId/socket", s.session.GetSocket)
	usersGroup.PUT( /*  */ "", s.session.PutUser, apimw.Session(s.dispatcher))

	pagesGroup := boardGroup.Group("/:id/pages", apimw.Session(s.dispatcher))
	pagesGroup.GET( /*  */ "", s.session.GetPageRank)
	pagesGroup.POST( /* */ "", s.session.PostPages)
	pagesGroup.PUT( /*  */ "", s.session.PutPages)
	pagesGroup.GET( /*  */ "/:pageId", s.session.GetPage)
	pagesGroup.GET( /*  */ "/sync", s.session.GetPageSync)
	pagesGroup.POST( /* */ "/sync", s.session.PostPageSync)

	attachGroup := boardGroup.Group("/:id/attachments", apimw.Session(s.dispatcher))
	attachGroup.POST( /**/ "", s.session.PostAttachment,
		libmw.RateLimiting(s.cfg.Server.RPM, libmw.WithUserIP()))
	attachGroup.GET( /* */ "/:attachId", s.session.GetAttachment)

	if s.cfg.Server.Metrics.Enabled {
		s.setMetricsRoutes()
	}

}

func (s *Server) setMetricsRoutes() {
	metricsGroup := s.echo.Group(s.cfg.Server.Metrics.Route,
		libmw.BasicAuth(s.cfg.Server.Metrics.User, s.cfg.Server.Metrics.Password))
	metricsGroup.GET("", s.metrics.GetMetrics)
}
