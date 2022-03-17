package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/heat1q/boardsite/api/middleware"
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b", echomw.Gzip(), middleware.RequestLogger())

	createGroup := boardGroup.Group("/create", middleware.RateLimiting(s.cfg.Server.RPM, middleware.WithIP()))
	createGroup.POST( /**/ "", s.session.PostCreateSession)

	configGroup := boardGroup.Group("/:id/config", middleware.Session(s.dispatcher))
	configGroup.GET( /*  */ "", s.session.GetSessionConfig)

	hostGroup := boardGroup.Group("",
		middleware.Session(s.dispatcher),
		middleware.Host(),
		middleware.GithubAuth(&s.cfg.Github, s.validator))
	hostGroup.PUT("/:id/config", s.session.PutSessionConfig)
	hostGroup.PUT("/:id/users/:userId", s.session.PutKickUser)

	usersGroup := boardGroup.Group("/:id/users")
	usersGroup.POST( /* */ "", s.session.PostUsers)
	usersGroup.GET( /*  */ "/:userId/socket", s.session.GetSocket)
	usersGroup.PUT( /*  */ "", s.session.PutUser, middleware.Session(s.dispatcher))

	pagesGroup := boardGroup.Group("/:id/pages", middleware.Session(s.dispatcher))
	pagesGroup.GET( /*  */ "", s.session.GetPageRank)
	pagesGroup.POST( /* */ "", s.session.PostPages)
	pagesGroup.PUT( /*  */ "", s.session.PutPages)
	pagesGroup.GET( /*  */ "/:pageId", s.session.GetPage)
	pagesGroup.GET( /*  */ "/sync", s.session.GetPageSync)
	pagesGroup.POST( /* */ "/sync", s.session.PostPageSync)

	attachGroup := boardGroup.Group("/:id/attachments", middleware.Session(s.dispatcher))
	attachGroup.POST( /**/ "", s.session.PostAttachment,
		middleware.GithubAuth(&s.cfg.Github, s.validator),
		middleware.RateLimiting(s.cfg.Server.RPM, middleware.WithUserIP()))
	attachGroup.GET( /* */ "/:attachId", s.session.GetAttachment)

	if s.cfg.Server.Metrics.Enabled {
		s.setMetricsRoutes()
	}

	if s.cfg.Github.Enabled {
		s.setGithubRoutes()
	}
}

func (s *Server) setMetricsRoutes() {
	metricsGroup := s.echo.Group(s.cfg.Server.Metrics.Route,
		middleware.BasicAuth(s.cfg.Server.Metrics.User, s.cfg.Server.Metrics.Password))
	metricsGroup.GET("", s.metrics.GetMetrics)
}

func (s *Server) setGithubRoutes() {
	githubGroup := s.echo.Group("/github/oauth", middleware.RequestLogger())
	githubGroup.GET("/authorize", s.github.GetAuthorize)
	githubGroup.GET("/callback", s.github.GetCallback)
	githubGroup.GET("/validate", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) }, middleware.GithubAuth(&s.cfg.Github, s.validator))
}
