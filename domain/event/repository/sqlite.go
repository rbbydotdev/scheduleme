package repository

import (
	"database/sql"
	"fmt"
	"scheduleme/domain/dbutil"
	"scheduleme/domain/event"
	"scheduleme/toerr"
	"scheduleme/values"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{
		db,
	}
}

type Event interface {
	CreateEvent(event *event.Event) (values.ID, error)
	AllForUserID(userID values.ID) (*event.Events, error)
	AttachRemoteByID(ID values.ID, ri *event.RouteInfo) error
	GetByID(id values.ID) (*event.Event, error)
	GetEventById(id values.ID) (*event.Event, error)
	GetEventByIDForUserID(id values.ID, userID values.ID) (*event.Event, error)
	UpdateEvent(event *event.Event) (int64, error)
	UpdateEventForUserID(event *event.Event, userID values.ID) (int64, error)
	DeleteEvent(id values.ID) (int64, error)
	DeleteEventForUserID(id values.ID, userID values.ID) (int64, error)
	AllPublicForUserID(userID values.ID) (*event.Events, error)
}

// CreateEvent inserts a new event into the database
func (r *repository) CreateEvent(event *event.Event) (values.ID, error) {
	res, err := r.db.Exec(`INSERT INTO events (name, duration, avail_masks, user_id, visible) VALUES (?, ?, ?, ?, ?)`,
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

func (r *repository) AllForUserID(userID values.ID) (*event.Events, error) {
	var events event.Events

	rows, err := r.db.Query(`SELECT id, name, duration, avail_masks, visible, user_id FROM events WHERE user_id = ?`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(fmt.Errorf("no events found for user_id=%v err=%w", userID, err)).Msg("no events found")
		}
		return nil, toerr.BadRequest(fmt.Errorf("error AllForUserID, db query: %w", err)).Msg("error getting events")
	}
	defer rows.Close()

	for rows.Next() {
		event := event.Event{}
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

// TODO: the route_info used here may indicate domain leakage
func (r *repository) AttachRemoteByID(ID values.ID, ri *route_info.RouteInfo) (err error) {
	e, err := r.GetByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(err)
		}
		return toerr.Internal(err)
	}
	ri.Event = *e
	return
}

func (r *repository) GetByID(id values.ID) (*event.Event, error) {
	e, err := r.GetEventById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(err)
		}
		return nil, toerr.Internal(err)
	}
	return e, nil
}

// GetEventById finds an event by ID
func (r *repository) GetEventById(id values.ID) (*event.Event, error) {
	event := &event.Event{}
	err := r.db.QueryRow(`SELECT id, name, duration, avail_masks, user_id, visible FROM events WHERE id = ?`, id).Scan(
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
func (r *repository) GetEventByIDForUserID(id values.ID, userID values.ID) (*event.Event, error) {
	event := &event.Event{}
	err := r.db.QueryRow(`SELECT id, name, duration, avail_masks, user_id, visible FROM events WHERE id = ? AND user_id = ?`, id, userID).Scan(
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
func (r *repository) UpdateEvent(event *event.Event) (int64, error) {
	return dbutil.WithCount(
		r.db.Exec(`UPDATE events SET name = ?, duration = ?, avail_masks = ?, visible = ? WHERE id = ?`,
			event.Name, event.Duration, event.AvailMasks, event.Visible, event.ID),
	)
}

// UpdateEventForUserID updates an existing event in the database for a given userID
func (r *repository) UpdateEventForUserID(event *event.Event, userID values.ID) (int64, error) {
	return dbutil.WithCount(
		r.db.Exec(`UPDATE events SET name = ?, duration = ?, avail_masks = ?, visible = ? WHERE id = ? AND user_id = ?`,
			event.Name, event.Duration, event.AvailMasks, event.Visible, event.ID, userID),
	)

}

// DeleteEvent deletes an event from the database
func (r *repository) DeleteEvent(id values.ID) (int64, error) {
	return dbutil.WithCount(
		r.db.Exec(`DELETE FROM events WHERE id = ?`, id),
	)
}

// DeleteEventForUserID deletes an event from the database for a given userID
func (r *repository) DeleteEventForUserID(id values.ID, userID values.ID) (int64, error) {
	return dbutil.WithCount(
		r.db.Exec(`DELETE FROM events WHERE id = ? AND user_id = ?`, id, userID),
	)
}

// AllPublicForUserID returns all public events for a given user
func (r *repository) AllPublicForUserID(userID values.ID) (*event.Events, error) {
	var events event.Events
	rows, err := r.db.Query(`SELECT id, name, duration, avail_masks, visible FROM events WHERE user_id = ? AND visible = ?`, userID, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := event.Event{}
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
