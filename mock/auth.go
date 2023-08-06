package mock

import (
	"time"

	"scheduleme/models"
	"scheduleme/values"

	"golang.org/x/oauth2"
)

// add mock interface
var _ models.AuthServiceInterface = (*AuthService)(nil)

type AuthService struct {
	CreateAuthFn                     func(auth *models.Auth) (values.ID, error)
	AttachRemoteByIDfn               func(ID values.OAuthSource, routeInfo *models.RouteInfo) error
	GetAuthByEmailFn                 func(email string) (*models.Auth, error)
	GetAuthByUserIDFn                func(userID values.ID, source values.OAuthSource) (*models.Auth, error)
	UpdateOrCreateAuthWithSourceIDFn func(userID values.ID, sourceID string, token *oauth2.Token, userInfo *models.UserInfo) (*models.Auth, error)
	GetAuthBySourceIDFn              func(sourceID string) (*models.Auth, error)
	UpdateAuthBySourceIDFn           func(sourceID string, accessToken values.Token, refreshToken values.Token, expiry time.Time) error
	UpdateAuthByIDFn                 func(ID values.ID, accessToken values.Token, refreshToken values.Token, expiry time.Time) error
}

func (a *AuthService) CreateAuth(auth *models.Auth) (values.ID, error) {
	return a.CreateAuthFn(auth)
}

func (a *AuthService) GetAuthByEmail(email string) (*models.Auth, error) {
	return a.GetAuthByEmailFn(email)
}

func (a *AuthService) GetAuthByUserID(authID values.ID, source values.OAuthSource) (*models.Auth, error) {
	return a.GetAuthByUserIDFn(authID, source)
}

func (a *AuthService) UpdateOrCreateAuthWithSourceID(userID values.ID, sourceID string, token *oauth2.Token, userInfo *models.UserInfo) (*models.Auth, error) {
	return a.UpdateOrCreateAuthWithSourceIDFn(userID, sourceID, token, userInfo)
}

func (a *AuthService) GetAuthBySourceID(sourceID string) (*models.Auth, error) {
	return a.GetAuthBySourceIDFn(sourceID)
}

func (a *AuthService) AttachRemoteByID(ID values.OAuthSource, routeInfo *models.RouteInfo) error {
	return a.AttachRemoteByIDfn(ID, routeInfo)
}

func (a *AuthService) UpdateAuthBySourceID(sourceID string, accessToken values.Token, refreshToken values.Token, expiry time.Time) error {
	return a.UpdateAuthBySourceIDFn(sourceID, accessToken, refreshToken, expiry)
}

func (a *AuthService) UpdateAuthByID(ID values.ID, accessToken values.Token, refreshToken values.Token,
	expiry time.Time) error {
	return a.UpdateAuthByIDFn(ID, accessToken, refreshToken, expiry)
}
