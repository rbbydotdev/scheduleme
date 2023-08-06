package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"scheduleme/config"
	"scheduleme/mock"
	"scheduleme/models"
	"scheduleme/secure_cookie"
	serve "scheduleme/services"
	_ "scheduleme/test"
	"scheduleme/toerr"
	"scheduleme/values"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

func TestRegisterEventRoutes(t *testing.T) {

	testUser := func(id values.ID) *models.User {
		return &models.User{
			ID:   id,
			Name: "testuser",
		}
	}
	testEvent := func(id values.ID) *models.Event {
		return &models.Event{
			ID:         id,
			Name:       "testevent",
			UserID:     values.ID(1),
			Duration:   1,
			AvailMasks: &values.AvailMasks{},
			Visible:    true,
		}
	}
	repo := &models.Repo{
		UserService: &mock.UserService{
			AttachRemoteByIDFn: func(id values.ID, ri *models.RouteInfo) (err error) {
				if id == 999 {
					return toerr.NotFound(nil)
				}
				u := testUser(id)
				ri.User = *u
				return
			},
		},
		EventService: &mock.EventService{
			EventsForUserRouteFn: func(ri *models.RouteInfo, ctx context.Context) error {
				return nil
			},
			AttachRemoteByIDFn: func(id values.ID, ri *models.RouteInfo) (err error) {
				if id == 777 {
					e := testEvent(id)
					e.Visible = false
					ri.Event = *e
					return
				}
				if id == 999 {
					return toerr.NotFound(nil)
				}
				e := testEvent(id)
				ri.Event = *e
				return
			},
		},
	}
	seccook := secure_cookie.NewSecureCookie("secret", "session")

	// srv := serve.BuildServer(seccook, repo, googleAuth)
	srv := serve.BuildServices(seccook, repo, &oauth2.Config{}, &config.ConfigStruct{})

	t.Run("Get Null User Events", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/events", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		eventsRoute := chi.NewRouter()
		eventsRoute.Route("/events", srv.RegisterEventsRoutes())
		eventsRoute.ServeHTTP(recorder, request)
		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if body := strings.TrimSpace(recorder.Body.String()); body != "null" {
			t.Errorf("handler returned unexpected body: got %v want %v", body, "null")
		}
	})

	t.Run("Get Events For User ID 1, StatusOK", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/users/1/events/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()

		router := chi.NewRouter()
		// router.Use(testhelpers.MockContextMiddleware(models.RouteInfo{User: models.User{ID: 1}, Event: models.Event{ID: 1}}))
		router.Route("/users", srv.RegisterUserRoutes())
		router.ServeHTTP(recorder, request)
		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		wantEvent := testEvent(values.ID(1))
		gotEvent := &models.Event{}
		err = json.NewDecoder(recorder.Body).Decode(gotEvent)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(gotEvent, wantEvent) {
			t.Errorf("handler returned unexpected event: got %+v, want %+v", gotEvent, wantEvent)
		}

	})

	t.Run("Get Private Event, StatusUnauthorized", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/users/777/events/777", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Route("/users", srv.RegisterUserRoutes())
		router.ServeHTTP(recorder, request)

		if got, want := recorder.Code, http.StatusUnauthorized; got != want {
			t.Errorf("handler returned wrong status code: got %v want %v", got, want)
		}
	})
	t.Run("Get Event, StatusNotFound", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/users/1/events/999", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Route("/users", srv.RegisterUserRoutes())
		router.ServeHTTP(recorder, request)

		if got, want := recorder.Code, http.StatusNotFound; got != want {
			t.Errorf("handler returned wrong status code: got %v want %v", got, want)
		}
	})

	t.Run("Get User, StatusNotFound", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/users/999/events/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Route("/users", srv.RegisterUserRoutes())
		router.ServeHTTP(recorder, request)

		if got, want := recorder.Code, http.StatusNotFound; got != want {
			t.Errorf("handler returned wrong status code: got %v want %v", got, want)
		}
	})
}
