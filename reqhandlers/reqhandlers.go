package reqhandlers

//TODO consolidate this package into another one?
import (
	"context"
	"encoding/json"
	"scheduleme/models"
	"scheduleme/toerr"
	"scheduleme/values"
	"strconv"

	"net/http"
	"net/url"
	"scheduleme/frame"

	"scheduleme/secure_cookie"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type Handlers struct {
	SecureCookie secure_cookie.SecureCookieInterface
}

// GoogleAuth   *oauth2.Config
func NewHandlers(sc secure_cookie.SecureCookieInterface, ga *oauth2.Config, r *models.Repo) *Handlers {
	return &Handlers{
		SecureCookie: sc,
	}
}

/*

This package contains handlers that are used to parse parameters and body data
This keeps request logistics out of the controllers, which deal more with routing the data
into and out of the models. This also allows for controllers to be more simple in that they
can simply pull out stored data from the context and pass it to the model methods. This allows
for some degree of DRYness.


TODO (possibly) add req validation to these handlers

*/

type ID = values.ID

type Queries[T any] interface {
	Parse(url.Values) error
	UpdatesQueryInfo(*models.QueryInfo, T)
	New() T
}

func ParseQuery[T Queries[T]]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var m T
			query := m.New()
			err := query.Parse(r.URL.Query())
			if err != nil {
				toerr.BadRequest(err).Render(w, r)
				return
			}
			ctx := frame.ModifyContextWith(r.Context(), func(qi *models.QueryInfo) {
				query.UpdatesQueryInfo(qi, query)
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type Mutatable interface {
	Validate() error
	ModifiesBodyInfo(bi *models.BodyInfo, ri models.RouteInfo, si models.SessionInfo)
}

// RouteInfoLoader is a middleware function that loads route information into the request context.
func RouteInfoLoader(fn func(*models.RouteInfo, context.Context) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newRi := frame.FromContext[models.RouteInfo](r.Context())
			err := fn(newRi, r.Context())
			if err != nil {
				toerr.Render(w, r, err).Log()
				return
			}
			frame.ServeWithNewContextInfo(w, r, next, newRi)
		})
	}
}

func ParseBody[T Mutatable]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var m T
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				toerr.Invalid(err).Msg("json decoding error, check document ;)").Render(w, r).Log()
				return
			}
			if err := m.Validate(); err != nil {
				toerr.Invalid(err).Msg("validation error").Msg(err.Error()).Render(w, r).Log()
				return
			}
			ri := frame.FromContext[models.RouteInfo](r.Context())
			si := frame.FromContext[models.SessionInfo](r.Context())
			frame.ModifyContextWith(r.Context(), func(bi *models.BodyInfo) {
				m.ModifiesBodyInfo(bi, *ri, *si)
			})
			next.ServeHTTP(w, r)
		})
	}
}

// Routable updates route info object, giving it a member of its own
// model after a successful look up by ID
// UserService, will essentially do: RouteInfo.User, err := UserService.GetUserByID(ID)
type Routable interface {
	AttachRemoteByID(values.ID, *models.RouteInfo) error
}

/*
Takes ID from url Path and looks up entity in service, then service will update route Info with this entity if its found
The following route handlers will then therefore have access to this entity
*/

func ResourceByID(paramName string, service Routable) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if pathID := chi.URLParam(r, paramName); pathID != "" {
				id, err := strconv.Atoi(pathID)
				pathIDtyped := values.ID(id)
				if err != nil {
					toerr.BadRequest(err).Render(w, r).Msg("invalid id, must be int").Log()
					return
				}
				newRi := frame.FromContext[models.RouteInfo](r.Context())
				err = service.AttachRemoteByID(pathIDtyped, newRi)
				if err != nil {
					toerr.Render(w, r, err).Log()
					return

				}
				frame.ServeWithNewContextInfo(w, r, next, newRi)

				//testnext.ServeHTTP(w, r.WithContext(NewContextWith(r.Context(), info)))

			}
		})
	}
}
