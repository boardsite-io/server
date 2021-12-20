package api

import "github.com/heat1q/boardsite/api/middleware"

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b", middleware.RequestLogger(s.echo.Logger))
	boardGroup.POST( /*  */ "/create", s.session.PostCreateSession)
	boardGroup.GET( /*   */ "/:id/users", s.session.GetUsers)
	boardGroup.POST( /*  */ "/:id/users", s.session.PostUsers)
	boardGroup.GET( /*   */ "/:id/users/:userId/socket", s.session.GetSocket)
	boardGroup.GET( /*   */ "/:id/pages", s.session.GetPages)
	boardGroup.POST( /*  */ "/:id/pages", s.session.PostPages)
	boardGroup.PUT( /*   */ "/:id/pages", s.session.PutPages)
	boardGroup.DELETE( /**/ "/:id/pages", s.session.DeletePages)

	boardGroup.GET( /*   */ "/:id/pages/:pageId", s.session.GetPageUpdate)
	boardGroup.DELETE( /**/ "/:id/pages/:pageId", s.session.DeletePageUpdate)

	boardGroup.POST( /**/ "/:id/attachments", s.session.PostAttachment)
	boardGroup.GET( /* */ "/:id/attachments/:attachId", s.session.GetAttachment)
}
