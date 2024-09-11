package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/secure_cookie"
	"scheduleme/toerr"
	"scheduleme/util"
)

type Middlewares struct {
	SecureCookie secure_cookie.SecureCookieInterface
}

func NewMiddleware(sc secure_cookie.SecureCookieInterface) *Middlewares {
	return &Middlewares{
		SecureCookie: sc,
	}
}

func (m *Middlewares) RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// reqi := frame.RequestInfoFromContext(r.Context())
		userID := models.UserIDFromContext(r.Context())
		if userID == 0 {
			//if request is type html set redirect url to current url and redirect to login
			//otherwise return unauthorized
			if util.IsHTML(r) {
				m.SecureCookie.PushRedirectURL(r, r.URL.Path)
				m.SecureCookie.HttpRedirect(w, r, "/login", http.StatusFound)
				return
			}
			toerr.Unauthorized(nil).Msg("unauthorized").Render(w, r)
			return

		} else {
			m.SecureCookie.NextServeHTTP(next, w, r)
		}

	})
}

// func ContentTypeNegotiationMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		allowedTypes := []string{"application/json", "text/html"}
// 		defaultType := "application/json"

// 		reqType := middleware.NegotiateContentType(r, allowedTypes, defaultType)

// 		if reqType != "application/json" && reqType != "text/html" {
// 			http.Error(w, "Invalid content type. Only 'application/json' and 'text/html' are accepted.", http.StatusBadRequest)
// 			return
// 		}

// 		// Pass the request along
// 		next.ServeHTTP(w, r)
// 	})
// }

func (m *Middlewares) PrintDebugSession(r *http.Request) {
	ses, _ := m.SecureCookie.GetSession(r)
	bytes, err := json.MarshalIndent(ses, "", "\t")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	log.Printf("%s\n", bytes)
}

//TODO flash and redirect url should probably be taken out of session stored in context in possibly another middleware function
//TODO cant figure out if data should be put into context and carried on or if

func (m *Middlewares) SessionMiddleware(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// userID, err := s.SecureCookie.GetSessionUserID(r)
			ses, err := m.SecureCookie.GetSession(r)
			if err != nil && err != http.ErrNoCookie {
				//if the session is weird, clear it and return
				//perhaps a redirect to login would be best here instead
				m.SecureCookie.ClearSession(w)
				toerr.Unauthorized(err).Msg("invalid session")
				return
			}
			if ses == nil {
				ses = models.NewSession()
			}
			sessionCtx := frame.NewContextWith(r.Context(), ses)
			next.ServeHTTP(w, r.WithContext(sessionCtx))

		})
	}
}
