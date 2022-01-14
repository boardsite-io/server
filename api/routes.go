package api

import "github.com/heat1q/boardsite/api/middleware"

const (
	rpmCreateBoard = 10
	rpmSession     = 100
)

// setRoutes sets the api routes
func (s *Server) setRoutes() {
	boardGroup := s.echo.Group("/b",
		middleware.Monitoring(),
		middleware.RequestLogger(),
		middleware.RateLimiting(rpmCreateBoard, middleware.WithIP()))
	boardGroup.POST( /*  */ "/create", s.session.PostCreateSession)
	boardGroup.POST( /*  */ "/:id/users", s.session.PostUsers)

	sessionGroup := s.echo.Group("/b/:id",
		middleware.Monitoring(),
		middleware.RequestLogger(),
		middleware.Session(s.dispatcher),
		middleware.RateLimiting(rpmSession, middleware.WithUserID()))
	sessionGroup.GET( /*   */ "/users", s.session.GetUsers)
	sessionGroup.GET( /*   */ "/users/:userId/socket", s.session.GetSocket)
	sessionGroup.GET( /*   */ "/pages", s.session.GetPages)
	sessionGroup.POST( /*  */ "/pages", s.session.PostPages)
	sessionGroup.PUT( /*   */ "/pages", s.session.PutPages)
	sessionGroup.DELETE( /**/ "/pages", s.session.DeletePages)

	sessionGroup.GET( /*   */ "/pages/:pageId", s.session.GetPageUpdate)
	sessionGroup.DELETE( /**/ "/pages/:pageId", s.session.DeletePageUpdate)

	sessionGroup.POST( /**/ "/attachments", s.session.PostAttachment)
	sessionGroup.GET( /* */ "/attachments/:attachId", s.session.GetAttachment)
}
