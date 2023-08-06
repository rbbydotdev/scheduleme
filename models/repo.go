package models

import (
	sq "scheduleme/sqlite"

	"golang.org/x/oauth2"
)

type Repo struct {
	UserService           UserServiceInterface
	EventService          EventServiceInterface
	AuthService           *AuthService
	GoogleCalendarService *GoogleCalendarService
	APIKeyService         *APIKeyService
	// PrivateEventService   *PrivateEventService
}

type GoogleUserInfo struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewRepo(db *sq.Db, oac *oauth2.Config) *Repo {
	return &Repo{
		UserService:           NewUserService(db),
		EventService:          NewEventService(db),
		AuthService:           NewAuthService(db), //TODO
		GoogleCalendarService: NewGoogleCalendarService(oac),
		APIKeyService:         NewAPIKeyService(db),
		// PrivateEventService:   NewPrivateEventService(db),
	}
}
