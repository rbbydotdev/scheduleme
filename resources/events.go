package resources

import (
	"fmt"
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/toerr"

	"github.com/go-chi/render"
)

func (re *Resources) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID := models.UserIDFromContext(r.Context())
	events, err := re.Repo.EventService.AllForUserID(userID)
	if err != nil {
		toerr.Render(w, r, err)
		return
	}
	render.JSON(w, r, events)
}

func (re *Resources) CreateEvent(w http.ResponseWriter, r *http.Request) {
	bi := frame.FromContext[models.BodyInfo](r.Context())
	_, err := re.Repo.EventService.CreateEvent(bi.Event)
	if err != nil {
		toerr.Internal(err).Msg("failed to create event").Render(w, r)
		return
	}
	render.Status(r, http.StatusCreated)
}

func (re *Resources) GetEvent(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.Event.View())
}
func (re *Resources) GetEvents(w http.ResponseWriter, r *http.Request) {
	var events models.EventsView
	ri := frame.FromContext[models.RouteInfo](r.Context())
	si := frame.FromContext[models.SessionInfo](r.Context())
	if ri.User.ID == si.UserID {
		events = ri.Events.ViewPrivate()
	} else {
		events = ri.Events.View()
	}
	render.JSON(w, r, events)
}
func (re *Resources) GetEventsPrivate(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.Events.ViewPrivate())
}

func (re *Resources) GetPublicEvents(w http.ResponseWriter, r *http.Request) {
	userID := models.UserIDFromContext(r.Context())
	events, err := re.Repo.EventService.AllForUserID(userID)
	if err != nil {
		toerr.Render(w, r, err)
		return
	}
	render.JSON(w, r, events)
}

func (re *Resources) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	bi := frame.FromContext[models.BodyInfo](r.Context())
	sessUserID := models.UserIDFromContext(r.Context())
	count, err := re.Repo.EventService.UpdateEventForUserID(bi.Event, sessUserID)
	if count == 0 {
		toerr.NotFound(fmt.Errorf("%v, event id %v for user id %v not found", err, sessUserID, bi.Event.ID)).Render(w, r)
		return
	}
	if err != nil {
		toerr.Internal(err).Msg("failed to update event").Render(w, r)
		return
	}
	render.JSON(w, r, &bi.Event)
}

func (re *Resources) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	sessUserID := models.UserIDFromContext(r.Context())
	_, err := re.Repo.EventService.DeleteEventForUserID(ri.Event.ID, sessUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render.NoContent(w, r)
}
