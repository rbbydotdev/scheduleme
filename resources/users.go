package resources

import (
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/values"

	"github.com/go-chi/render"
)

func (re *Resources) GetUserPublic(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.User.View())
}
func (re *Resources) GetUserPrivate(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	render.JSON(w, r, ri.User.ViewPrivate())
}

func (re *Resources) UpdateUser(w http.ResponseWriter, r *http.Request) {
	bi := frame.FromContext[models.BodyInfo](r.Context())
	count, err := re.Repo.UserService.UpdateUser(bi.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	render.JSON(w, r, &bi.User)
}

func (re *Resources) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	count, err := re.Repo.UserService.DeleteUser(values.ID(ri.User.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	render.NoContent(w, r)
}
