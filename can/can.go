package can

import (
	"context"
	"net/http"
	"scheduleme/frame"
	"scheduleme/hof"
	"scheduleme/models"
	"scheduleme/toerr"
)

func mutateAPIKey(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.User.ID == sessUserID
}

func deleteUser(ctx context.Context) bool {
	// user := frame.UserFromContext(ctx)
	return frame.FromContext[models.SessionInfo](ctx).IsAdmin
	// return models.SessionInfoCtx(ctx).IsAdmin
}
func mutateUser(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.User.ID == sessUserID
}

func showEvents(ctx context.Context) bool {
	return true
}

func showEventsPrivate(ctx context.Context) bool {
	ri := frame.FromContext[models.RouteInfo](ctx)
	sessUserID := models.UserIDFromContext(ctx)
	return hof.Any(ri.Events, func(e models.Event) bool {
		return e.UserID != sessUserID
	})
}

func showUser(ctx context.Context) bool {
	return true
}
func showUserPrivate(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.User.ID == sessUserID

}

func mutateEvent(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.Event.UserID == sessUserID
}

func deleteEvent(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.Event.UserID == sessUserID
}

func showEvent(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	ri := frame.FromContext[models.RouteInfo](ctx)
	return ri.Event.Visible || ri.Event.UserID == sessUserID
}

func createEvent(ctx context.Context) bool {
	sessUserID := models.UserIDFromContext(ctx)
	return sessUserID != 0
}
func CanHandler(canFn func(context.Context) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !canFn(r.Context()) {
				toerr.Unauthorized(nil).Render(w, r).Log()
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

var ShowEvent = CanHandler(showEvent)
var DeleteEvent = CanHandler(deleteEvent)
var MutateEvent = CanHandler(mutateEvent)
var MutateUser = CanHandler(mutateUser)
var ShowUser = CanHandler(showUser)
var ShowPrivateUser = CanHandler(showUserPrivate)

var ShowAllEvents = CanHandler(showEventsPrivate)
var CreateEvent = CanHandler(createEvent)
var ShowEvents = CanHandler(showEvents)
var DeleteUser = CanHandler(deleteUser)

var MutateAPIKey = CanHandler(mutateAPIKey)
