package event

import (
	"scheduleme/domain/body_info"
	"scheduleme/domain/route_info"
	"scheduleme/domain/session_info"
	sq "scheduleme/sqlite"
	"scheduleme/values"
	"time"
)

type EventWithAvailability struct {
	Event        Event            `json:"event"`
	Availability values.DateSlots `json:"availability"`
}

func NewEventWithAvailability(event Event, availability values.DateSlots) *EventWithAvailability {
	return &EventWithAvailability{Event: event, Availability: availability}
}

// TODO: interesting usage of 'views' but this should really be done in the db query
func (es *Events) View() EventsView {
	var eventsView EventsView
	for i := 0; i < len(*es); i++ {
		event := (*es)[i]
		if event.Visible {
			eventsView = append(eventsView, event.View())
		}
	}
	return eventsView
}

func (es *Events) ViewPrivate() EventsView {
	var eventsView EventsView
	for i := 0; i < len(*es); i++ {
		event := (*es)[i]
		eventsView = append(eventsView, event.View())
	}
	return eventsView
}

func (e *Event) View() EventView {
	return EventView{
		ID:         e.ID,
		Name:       e.Name,
		Duration:   e.Duration,
		AvailMasks: e.AvailMasks,
		UserID:     e.UserID,
		Visible:    e.Visible,
	}
}

type Event struct {
	ID         values.ID          `json:"id"`
	Name       string             `json:"name"`
	Duration   time.Duration      `json:"duration"`
	AvailMasks *values.AvailMasks `json:"avail_masks"`
	UserID     values.ID          `json:"user_id"`
	Visible    bool               `json:"visible"`
}

type EventsView []EventView

type Events []Event

type EventService struct {
	db *sq.Db
}

type EventCreate struct {
	Name       string            `json:"name"`
	Duration   time.Duration     `json:"duration"`
	AvailMasks values.AvailMasks `json:"avail_masks"`
	Visible    bool              `json:"visible"`
}

type EventView struct {
	ID         values.ID          `json:"id"`
	Name       string             `json:"name"`
	Duration   time.Duration      `json:"duration"`
	AvailMasks *values.AvailMasks `json:"avail_masks"`
	UserID     values.ID          `json:"user_id"`
	Visible    bool               `json:"visible"`
}

type EventMutate struct {
	Name       string             `json:"name"`
	Duration   time.Duration      `json:"duration"`
	AvailMasks *values.AvailMasks `json:"avail_masks"`
	Visible    bool               `json:"visible"`
}

func (em *EventMutate) ModifiesBodyInfo(bi *body_info.BodyInfo, ri *route_info.RouteInfo, si *session_info.SessionInfo) {
	bi.Event = &Event{
		ID:         ri.Event.ID,
		Name:       em.Name,
		Duration:   em.Duration,
		AvailMasks: em.AvailMasks,
		Visible:    em.Visible,
		UserID:     si.UserID,
	}
}

func (ec *EventCreate) Validate() error {
	return nil
}

func (ec *EventCreate) ModifiesBodyInfo(bi *body_info.BodyInfo, ri *route_info.RouteInfo, si *session_info.SessionInfo) {
	bi.Event = &Event{
		Name:       ec.Name,
		Duration:   ec.Duration,
		AvailMasks: &ec.AvailMasks,
		Visible:    ec.Visible,
		UserID:     si.UserID,
	}
}

type RouteInfo struct {
	Event  Event
	Events Events
	Offset int
	Page   int
	Filter string
}

func (ri RouteInfo) ContextKey() string {
	return "EventRouteInfo"
}

func NewRouteInfo() *RouteInfo {
	return &RouteInfo{}
}
