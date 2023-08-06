package services

import (
	"net/http"
	"scheduleme/secure_cookie"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func TopRoutes(topSvc *Services) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(topSvc.Middlewares.SessionMiddleware(topSvc.ServiceConfig.Secret))
	// r.Use(frame.RouteInfoMiddleware)
	// r.Use(middlewares.AcceptMiddleware) //TODO: not needed probably for now

	r.Group(func(r chi.Router) {
		//HTML PAGES
		r.Use(render.SetContentType(render.ContentTypeHTML))
		r.Group(topSvc.RegisterOAuthRouters())
		r.Group(topSvc.RegisterSessionRouters())
		r.Group(RegisterPageRouters())
		if topSvc.ServiceConfig.ENV.IsDev() {
			r.Route("/debug", RegisterDebugRoutes(topSvc.SecureCookie.(*secure_cookie.SecureCookie)))
		}
	})
	r.Route("/api", func(r chi.Router) {
		//JSON API
		r.Use(render.SetContentType(render.ContentTypeJSON))
		// r.Get("/ping", topSvc.HandlePing)

		r.Route("/users",
			topSvc.RegisterUserRoutes(),
		)
		r.Route("/events",
			topSvc.RegisterEventsRoutes(),
		)
		r.Route("/apikeys",
			topSvc.RegisterApiKeyRoutes(),
		)

	})
	return r
}
