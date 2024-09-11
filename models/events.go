package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"scheduleme/frame"
	sq "scheduleme/sqlite"
	"scheduleme/toerr"
	"scheduleme/values"
	"time"
)

func NewEventService(db *sq.Db) EventServiceInterface {
	return &EventService{db: db}
}

func (e *Event) Validate() error {
	//validate duration
	//duration must only be 15,30,45,60,75,90,105,120,135,150,165,180
	if e.Duration%(15*time.Minute) != 0 {
		return errors.New("invalid duration")
	}
	if e.Duration < 15*time.Minute || e.Duration > 180*time.Minute {
		return errors.New("invalid duration")
	}
	return nil
}

//TODO use AvalMask to determine if event is available

type Events []Event
type EventsView []EventView

type EventWithAvailability struct {
	Event        Event            `json:"event"`
	Availability values.DateSlots `json:"availability"`
}

func NewEventWithAvailability(event Event, availability values.DateSlots) *EventWithAvailability {
	return &EventWithAvailability{Event: event, Availability: availability}
}

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

func (em *EventMutate) Validate() error {
	return nil
}

func (em *EventMutate) ModifiesBodyInfo(bi *BodyInfo, ri RouteInfo, si SessionInfo) {
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

func (ec *EventCreate) ModifiesBodyInfo(bi *BodyInfo, ri RouteInfo, si SessionInfo) {
	bi.Event = &Event{
		Name:       ec.Name,
		Duration:   ec.Duration,
		AvailMasks: &ec.AvailMasks,
		Visible:    ec.Visible,
		UserID:     si.UserID,
	}
}

type EventCreate struct {
	Name       string            `json:"name"`
	Duration   time.Duration     `json:"duration"`
	AvailMasks values.AvailMasks `json:"avail_masks"`
	Visible    bool              `json:"visible"`
}

type EventService struct {
	db *sq.Db
}

type EventServiceInterface interface {
	AllForUserID(userID values.ID) (*Events, error)
	CreateEvent(event *Event) (values.ID, error)
	GetEventById(id values.ID) (*Event, error)
	UpdateEvent(event *Event) (int64, error)
	UpdateEventForUserID(event *Event, userID values.ID) (int64, error)
	DeleteEvent(id values.ID) (int64, error)
	DeleteEventForUserID(id values.ID, userID values.ID) (int64, error)
	AllPublicForUserID(userID values.ID) (*Events, error)
	EventsForUserRoute(ri *RouteInfo, ctx context.Context) error
	AttachRemoteByID(ID values.ID, ri *RouteInfo) error
}

func (es *EventService) EventsForUserRoute(ri *RouteInfo, ctx context.Context) (err error) {
	sessUserID := frame.FromContext[SessionInfo](ctx).UserID
	var events *Events
	if ri.User.ID == sessUserID && ri.User.ID != 0 {
		events, err = es.AllForUserID(ri.User.ID)
	} else {
		events, err = es.AllPublicForUserID(ri.User.ID)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			toerr.NotFound(err).Msg("no events found")
		}
		return toerr.Internal(err).Msg("error getting events")
	}
	//print user id of each event
	for i := 0; i < len(*events); i++ {
		log.Printf("event %v user id: %v\n", (*events)[i].ID, (*events)[i].UserID)
	}
	ri.Events = *events
	return
}

// CreateEvent inserts a new event into the database
func (s *EventService) CreateEvent(event *Event) (values.ID, error) {
	res, err := s.db.Exec(`INSERT INTO events (name, duration, avail_masks, user_id, visible) VALUES (?, ?, ?, ?, ?)`,
		event.Name, event.Duration, event.AvailMasks, event.UserID, event.Visible)
	if err != nil {
		return 0, toerr.Internal(fmt.Errorf("failed to create event err=%w", err))
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, toerr.Internal(fmt.Errorf("failed to create event err=%w", err))
	}
	return values.ID(id), nil
}

func (s *EventService) AllForUserID(userID values.ID) (*Events, error) {
	var events Events

	rows, err := s.db.Query(`SELECT id, name, duration, avail_masks, visible, user_id FROM events WHERE user_id = ?`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(fmt.Errorf("no events found for user_id=%v err=%w", userID, err)).Msg("no events found")
		}
		return nil, toerr.BadRequest(fmt.Errorf("error AllForUserID, db query: %w", err)).Msg("error getting events")
	}
	defer rows.Close()

	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.ID, &event.Name, &event.Duration, &event.AvailMasks, &event.Visible, &event.UserID)
		if err != nil {
			return nil, toerr.Internal(err)
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, toerr.Internal(err)
	}
	return &events, nil
}

func (es *EventService) AttachRemoteByID(ID values.ID, ri *RouteInfo) (err error) {
	e, err := es.GetByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(err)
		}
		return toerr.Internal(err)
	}
	ri.Event = *e
	return
}

func (es *EventService) GetByID(id values.ID) (*Event, error) {
	e, err := es.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(err)
		}
		return nil, toerr.Internal(err)
	}
	return e, nil
}

// GetEventById finds an event by ID
func (es *EventService) GetEventById(id values.ID) (*Event, error) {
	event := &Event{}
	err := es.db.QueryRow(`SELECT id, name, duration, avail_masks, user_id, visible FROM events WHERE id = ?`, id).Scan(
		&event.ID, &event.Name, &event.Duration, &event.AvailMasks, &event.UserID, &event.Visible)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(err)
		}
		return nil, toerr.Internal(err)
	}
	return event, nil
}

// GetEventByIDForUSerID finds an event by ID for UserID
func (es *EventService) GetEventByIDForUserID(id values.ID, userID values.ID) (*Event, error) {
	event := &Event{}
	err := es.db.QueryRow(`SELECT id, name, duration, avail_masks, user_id, visible FROM events WHERE id = ? AND user_id = ?`, id, userID).Scan(
		&event.ID, &event.Name, &event.Duration, &event.AvailMasks, &event.UserID, &event.Visible)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(err).Msg("event id %v for user id %v not found", id, userID) //Should this just say 'not found' ? or be more specific?
		}
		return nil, toerr.Internal(err)
	}
	return event, nil
}

// UpdateEvent updates an existing event in the database
func (es *EventService) UpdateEvent(event *Event) (int64, error) {
	return withCount(
		es.db.Exec(`UPDATE events SET name = ?, duration = ?, avail_masks = ?, visible = ? WHERE id = ?`,
			event.Name, event.Duration, event.AvailMasks, event.Visible, event.ID),
	)
}

// UpdateEventForUserID updates an existing event in the database for a given userID
func (es *EventService) UpdateEventForUserID(event *Event, userID values.ID) (int64, error) {
	return withCount(
		es.db.Exec(`UPDATE events SET name = ?, duration = ?, avail_masks = ?, visible = ? WHERE id = ? AND user_id = ?`,
			event.Name, event.Duration, event.AvailMasks, event.Visible, event.ID, userID),
	)

}

// DeleteEvent deletes an event from the database
func (es *EventService) DeleteEvent(id values.ID) (int64, error) {
	return withCount(
		es.db.Exec(`DELETE FROM events WHERE id = ?`, id),
	)
}

// DeleteEventForUserID deletes an event from the database for a given userID
func (es *EventService) DeleteEventForUserID(id values.ID, userID values.ID) (int64, error) {
	return withCount(
		es.db.Exec(`DELETE FROM events WHERE id = ? AND user_id = ?`, id, userID),
	)
}

// AllPublicForUserID returns all public events for a given user
func (es *EventService) AllPublicForUserID(userID values.ID) (*Events, error) {
	var events Events
	rows, err := es.db.Query(`SELECT id, name, duration, avail_masks, visible FROM events WHERE user_id = ? AND visible = ?`, userID, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.ID, &event.Name, &event.Duration, &event.AvailMasks, &event.Visible)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &events, nil
}
