package resources

import (
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"

	"github.com/go-chi/render"
)

func (re *Resources) ListAvailability(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.Availability)
}

func (re *Resources) GetEventAndAvailabilities(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, models.NewEventWithAvailability(ri.Event, ri.Availability))
}
