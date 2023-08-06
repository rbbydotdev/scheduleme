package resources

import (
	"scheduleme/models"
)

type Resources struct {
	Repo *models.Repo
}

func NewResources(r *models.Repo) *Resources {
	return &Resources{r}
}
