package api

import "github.com/heat1q/boardsite/api/middleware"

const (
	rpmSession = 100
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b", middleware.Monitoring(), middleware.RequestLogger())
	boardGroup.POST( /*  */ "/create", s.session.PostCreateSession, middleware.RateLimiting(1, middleware.WithIP()))

	usersGroup := boardGroup.Group("/:id/users", middleware.RateLimiting(rpmSession, middleware.WithUserIP()))
	usersGroup.POST( /*  */ "", s.session.PostUsers)
	usersGroup.GET( /*   */ "", s.session.GetUsers, middleware.Session(s.dispatcher))
	usersGroup.GET( /*   */ "/:userId/socket", s.session.GetSocket, middleware.Session(s.dispatcher))

	pagesGroup := boardGroup.Group("/:id/pages", middleware.Session(s.dispatcher), middleware.RateLimiting(rpmSession, middleware.WithUserID()))
	pagesGroup.GET( /*   */ "", s.session.GetPages)
	pagesGroup.POST( /*  */ "", s.session.PostPages)
	pagesGroup.PUT( /*   */ "", s.session.PutPages)
	pagesGroup.DELETE( /**/ "", s.session.DeletePages)
	pagesGroup.GET( /*   */ "/:pageId", s.session.GetPageUpdate)
	pagesGroup.DELETE( /**/ "/:pageId", s.session.DeletePageUpdate)

	attachGroup := boardGroup.Group("/:id/attachments", middleware.Session(s.dispatcher))
	attachGroup.POST( /**/ "", s.session.PostAttachment, middleware.RateLimiting(1, middleware.WithUserID()))
	attachGroup.GET( /* */ "/:attachId", s.session.GetAttachment)
}
