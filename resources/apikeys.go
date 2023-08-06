package resources

import (
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/toerr"

	"github.com/go-chi/render"
)

func (re *Resources) CreateApiKey(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	userID := ri.User.ID
	apiKey, err := re.Repo.APIKeyService.CreateAPIKeyForUserID(userID)
	if err != nil {
		toerr.Internal(err).Msg("failed to create api key").Render(w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &apiKey)
}

func (re *Resources) DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	apiKey, err := re.Repo.APIKeyService.DeleteAPIKeyByID(ri.APIKey.ID)
	if err != nil {
		toerr.Internal(err).Msg("failed to delete api key").Render(w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &apiKey)
}
