package mock

import (
	"scheduleme/models"

	"golang.org/x/oauth2"
)

type AuthUserFlow struct {
	AuthUserFn func(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error)
}

func (auf *AuthUserFlow) AuthUser(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error) {
	return auf.AuthUserFn(userInfo, token)
}
