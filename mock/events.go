package mock

import (
	// "database/sql"
	"context"
	"scheduleme/models"
	"scheduleme/values"
)

var _ models.EventServiceInterface = (*EventService)(nil)

type EventService struct {
	AllFn                  func() ([]*models.Event, error)
	AllForUserIDFn         func(userID values.ID) (*models.Events, error)
	GetEventByIdFn         func(id values.ID) (*models.Event, error)
	GetAllPublicEventsFn   func(userID values.ID) ([]*models.Event, error)
	UpdateEventForUserIDFn func(event *models.Event, userID values.ID) (int64, error)
	DeleteEventFn          func(id values.ID) (int64, error)
	DeleteEventForUserIDFn func(id values.ID, userID values.ID) (int64, error)
	AllPublicForUserIDFn   func(userID values.ID) (*models.Events, error)
	AttachRemoteByIDFn     func(id values.ID, ri *models.RouteInfo) error
	EventsForUserRouteFn   func(ri *models.RouteInfo, ctx context.Context) error
	UpdateEventFn          func(event *models.Event) (int64, error)
	CreateEventFn          func(event *models.Event) (values.ID, error)
}

func (s *EventService) EventsForUserRoute(ri *models.RouteInfo, ctx context.Context) error {
	return s.EventsForUserRouteFn(ri, ctx)
}
func (s *EventService) AttachRemoteByID(id values.ID, ri *models.RouteInfo) error {
	return s.AttachRemoteByIDFn(id, ri)
}

func (s *EventService) AllPublicForUserID(userID values.ID) (*models.Events, error) {
	return s.AllPublicForUserIDFn(userID)
}

func (s *EventService) All() ([]*models.Event, error) {
	return s.AllFn()
}

func (s *EventService) AllForUserID(userID values.ID) (*models.Events, error) {
	return s.AllForUserIDFn(userID)
}

func (s *EventService) GetAllPublicEvents(userID values.ID) ([]*models.Event, error) {
	return s.GetAllPublicEventsFn(userID)
}

func (s *EventService) CreateEvent(event *models.Event) (values.ID, error) {
	return s.CreateEventFn(event)
}

func (s *EventService) GetEventById(id values.ID) (*models.Event, error) {
	return s.GetEventByIdFn(id)
}

func (s *EventService) UpdateEvent(event *models.Event) (int64, error) {
	return s.UpdateEventFn(event)
}

func (s *EventService) UpdateEventForUserID(event *models.Event, userID values.ID) (int64, error) {
	return s.UpdateEventForUserIDFn(event, userID)
}

func (s *EventService) DeleteEvent(id values.ID) (int64, error) {
	return s.DeleteEventFn(id)
}

func (s *EventService) DeleteEventForUserID(id values.ID, userID values.ID) (int64, error) {
	return s.DeleteEventForUserIDFn(id, userID)
}
