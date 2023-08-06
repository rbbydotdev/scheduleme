package services

import (
	"scheduleme/secure_cookie"

	"github.com/go-chi/chi/v5"
)

func RegisterDebugRoutes(secureCookie *secure_cookie.SecureCookie) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/session", HandleSessionDebug(secureCookie))
		r.Get("/context", HandleContextDebug)
		r.Get("/sessioncontext", HandleSessionContextDebug)
	}

}
func RegisterPageRouters() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/home", HandleHome)
		r.Get("/ping", HandlePing)
		r.Get("/", RedirectHome)
	}
}
