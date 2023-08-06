package services

import (
	"scheduleme/can"
	"scheduleme/models"
	"scheduleme/reqhandlers"

	"github.com/go-chi/chi/v5"
)

func (srv *Services) RegisterOAuthRouters() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get(srv.ServiceConfig.GoogleRedirectPath, srv.OAuth.HandleGoogleCallback(srv.OAuth.HandleLoginSuccess))
	}
}
func (srv *Services) RegisterSessionRouters() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/login", srv.OAuth.HandleGoogleLogin)
		r.Get("/logout", srv.OAuth.HandleLogout)
	}
}
func (srv *Services) RegisterApiKeyRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(
				srv.Middlewares.RequireAuthMiddleware,
				reqhandlers.RouteInfoLoader(srv.Resources.Repo.UserService.MeForUserRoute),
				can.MutateAPIKey,
			)
			r.Post("/", srv.Resources.CreateApiKey)
			r.Route("/{ApiKeyID}", func(r chi.Router) {
				r.With(
					reqhandlers.ResourceByID("ApiKeyID", srv.Resources.Repo.APIKeyService),
				).Delete("/", srv.Resources.DeleteApiKey)
			})
		})
	}
}

func (srv *Services) RegisterUserRoutes() func(chi.Router) {
	return func(r chi.Router) {

		r.Route("/me", func(r chi.Router) {
			r.Use(
				srv.Middlewares.RequireAuthMiddleware,
				reqhandlers.RouteInfoLoader(srv.Resources.Repo.UserService.MeForUserRoute),
				can.ShowPrivateUser,
			)
			r.Get("/", srv.Resources.GetUserPrivate)
			r.Route("/events", srv.RegisterEventsRoutes())
		})

		r.Route("/{UserID}", func(r chi.Router) {
			r.Use(
				reqhandlers.ResourceByID("UserID", srv.Resources.Repo.UserService),
			)
			r.Get("/", srv.Resources.GetUserPublic)
			r.Route("/events",
				srv.RegisterEventsRoutes(),
			)
			r.With(
				srv.Middlewares.RequireAuthMiddleware,
				reqhandlers.ParseBody[*models.UserMutate](),
				can.MutateUser,
			).Put("/update", srv.Resources.UpdateUser)

			r.With(
				srv.Middlewares.RequireAuthMiddleware,
				can.DeleteUser,
			).Delete("/delete", srv.Resources.DeleteUser)
		})
	}
}

func (srv *Services) RegisterEventsRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.With(
			reqhandlers.RouteInfoLoader(srv.Resources.Repo.EventService.EventsForUserRoute),
			can.ShowEvents,
		).Get("/", srv.Resources.GetEvents)

		r.Route("/{EventID}", func(r chi.Router) {
			r.Use(reqhandlers.ResourceByID("EventID", srv.Resources.Repo.EventService))
			r.Route("/avail", func(r chi.Router) {
				r.With(
					reqhandlers.ParseQuery[*models.AvailQuery](),
					reqhandlers.RouteInfoLoader(srv.Resources.Repo.AuthService.AuthForUserRoute),
					reqhandlers.RouteInfoLoader(srv.Resources.Repo.GoogleCalendarService.AvailabilityForEventRoute),
					can.ShowEvent,
				).Get("/", srv.Resources.GetEventAndAvailabilities)
			})
			r.With(
				srv.Middlewares.RequireAuthMiddleware,
				reqhandlers.ParseBody[*models.EventMutate](),
				can.MutateEvent,
			).Put("/", srv.Resources.UpdateEvent)
			r.With(
				srv.Middlewares.RequireAuthMiddleware,
				can.DeleteEvent,
			).Delete("/", srv.Resources.DeleteEvent)
			r.With(
				can.ShowEvent,
			).Get("/", srv.Resources.GetEvent)
		})
	}
}
