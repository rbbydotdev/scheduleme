package handler

import (
	"fmt"
	"net/http"
	"scheduleme/domain/event"
	"scheduleme/domain/event/use_case"
	"scheduleme/domain/session_info"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/toerr"

	"github.com/go-chi/render"
)

type Handler struct {
	useCase use_case.Event
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID := models.UserIDFromContext(r.Context())
	events, err := h.useCase.AllForUserID(userID)
	if err != nil {
		toerr.Render(w, r, err)
		return
	}
	render.JSON(w, r, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	bi := frame.FromContext[event.BodyInfo](r.Context())
	_, err := h.useCase.CreateEvent(bi.Event)
	if err != nil {
		toerr.Internal(err).Msg("failed to create event").Render(w, r)
		return
	}
	render.Status(r, http.StatusCreated)
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.Event.View())
}
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	var events event.EventsView
	ri := frame.FromContext[event.RouteInfo](r.Context())
	si := frame.FromContext[session_info.SessionInfo](r.Context())
	if ri.User.ID == si.UserID {
		events = ri.Events.ViewPrivate()

	} else {
		events = ri.Events.View()
	}
	//render.JSON
	render.JSON(w, r, events)
}
func (h *Handler) GetEventsPrivate(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.Events.ViewPrivate())
}

func (h *Handler) GetPublicEvents(w http.ResponseWriter, r *http.Request) {
	userID := models.UserIDFromContext(r.Context()) //TODO
	events, err := h.useCase.AllForUserID(userID)
	if err != nil {
		toerr.Render(w, r, err)
		return
	}
	render.JSON(w, r, events)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	bi := frame.FromContext[models.BodyInfo](r.Context())
	sessUserID := models.UserIDFromContext(r.Context()) //TODO
	count, err := h.useCase.UpdateEventForUserID(bi.Event, sessUserID)
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

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	sessUserID := models.UserIDFromContext(r.Context())
	_, err := h.useCase.DeleteEventForUserID(ri.Event.ID, sessUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	render.NoContent(w, r)
}
