package use_case

import (
	"scheduleme/domain/event"
	"scheduleme/domain/event/repository"
	"scheduleme/values"
)

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

type EventUseCase struct {
	repo repository.Event
}

func (euc *EventUseCase) CreateEvent(event *event.Event) (values.ID, error) {
	return euc.repo.CreateEvent(event)
}
func (euc *EventUseCase) AllForUserID(userID values.ID) (*event.Events, error) {
	return euc.repo.AllForUserID(userID)
}
func (euc *EventUseCase) AttachRemoteByID(ID values.ID, ri *event.RouteInfo) error {
	return euc.repo.AttachRemoteByID(ID, ri)
}
func (euc *EventUseCase) GetByID(id values.ID) (*event.Event, error) {
	return euc.repo.GetByID(id)
}
func (euc *EventUseCase) GetEventById(id values.ID) (*event.Event, error) {
	return euc.repo.GetEventById(id)
}
func (euc *EventUseCase) GetEventByIDForUserID(id values.ID, userID values.ID) (*event.Event, error) {
	return euc.repo.GetEventByIDForUserID(id, userID)
}
func (euc *EventUseCase) UpdateEvent(event *event.Event) (int64, error) {
	return euc.repo.UpdateEvent(event)
}
func (euc *EventUseCase) UpdateEventForUserID(event *event.Event, userID values.ID) (int64, error) {
	return euc.repo.UpdateEventForUserID(event, userID)
}
func (euc *EventUseCase) DeleteEvent(id values.ID) (int64, error) {
	return euc.repo.DeleteEvent(id)
}
func (euc *EventUseCase) DeleteEventForUserID(id values.ID, userID values.ID) (int64, error) {
	return euc.repo.DeleteEventForUserID(id, userID)
}
func (euc *EventUseCase) AllPublicForUserID(userID values.ID) (*event.Events, error) {
	return euc.repo.AllPublicForUserID(userID)
}
