package services_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"scheduleme/config"
	"scheduleme/mock"
	"scheduleme/models"
	serve "scheduleme/services"
	"scheduleme/values"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRegisterUserRoutes(t *testing.T) {
	repo := &models.Repo{
		UserService: &mock.UserService{
			AttachRemoteByIDFn: func(id values.ID, ri *models.RouteInfo) (err error) {
				u := &models.User{
					ID:   id,
					Name: "test",
				}
				ri.User = *u
				return
			},
		},

		EventService: &mock.EventService{
			EventsForUserRouteFn: func(ri *models.RouteInfo, ctx context.Context) error {
				return nil
			},
		},
	}
	srv := serve.BuildServices(nil, repo, nil, &config.ConfigStruct{})

	t.Run("Get User 1", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		eventsRoute := chi.NewRouter()
		eventsRoute.Route("/users", srv.RegisterUserRoutes())
		eventsRoute.ServeHTTP(recorder, request)
		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if got, want := strings.TrimSpace(recorder.Body.String()), `{"id":1,"name":"test"}`; got != want {
			t.Errorf("handler returned unexpected body: got %v want %v", got, want)

		}
	})

}
