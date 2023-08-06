package models_test

import (
	"reflect"
	"scheduleme/models"
	"scheduleme/sqlite"
	"scheduleme/test"
	"scheduleme/values"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCreateGetEventByID(t *testing.T) {
	db, err := sqlite.NewOpenDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	u := models.NewUserService(db)
	uid, err := u.CreateUser(&models.User{
		Name:  "Test User",
		Email: "email@email.com",
	})
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	es := models.NewEventService(db)

	mask := test.MakeMask()
	event := &models.Event{
		Name:       "Test Event",
		Duration:   time.Hour,
		AvailMasks: &values.AvailMasks{mask},
		UserID:     uid,
		Visible:    true,
	}

	id, err := es.CreateEvent(event)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	event.ID = id

	got, err := es.GetEventById(id)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if !reflect.DeepEqual(got, event) {
		t.Errorf("Expected event %v, but got %v", event, got)
	}
	if !got.Visible {
		t.Errorf("Expected event to be visible, but got %v", got)
	}
}

func TestCreateEvent(t *testing.T) {

	db, err := sqlite.NewOpenDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	u := models.NewUserService(db)
	uid, _ := u.CreateUser(&models.User{
		Name:  "Test User",
		Email: "email@email.com",
	})

	s := models.NewEventService(db)
	mask := test.MakeMask()

	event := &models.Event{
		Name:       "Test Event",
		Duration:   time.Hour,
		AvailMasks: &values.AvailMasks{mask},
		UserID:     uid,
	}

	id, err := s.CreateEvent(event)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}
	event.ID = id

	got, err := s.GetEventById(1)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	want := event
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CreateEvent() mismatch (-want +got):\n%s", diff)
	}

}
