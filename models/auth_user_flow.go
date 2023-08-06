package models

import (
	"fmt"

	"golang.org/x/oauth2"
)

type AuthUserFlow struct {
	AuthService AuthServiceInterface
	UserService UserServiceInterface
}

func NewAuthUserFlow(
	authService AuthServiceInterface,
	userService UserServiceInterface,
) *AuthUserFlow {
	return &AuthUserFlow{
		AuthService: authService,
		UserService: userService,
	}
}

var _ AuthUserFlowInterface = (*AuthUserFlow)(nil)

// var GoogleAgent AgentInterface = (*GoogleAgent)(nil)

type AuthUserFlowInterface interface {
	AuthUser(userInfo *UserInfo, token *oauth2.Token) (*Auth, *User, error)
}

func (oa *AuthUserFlow) AuthUser(userInfo *UserInfo, token *oauth2.Token) (*Auth, *User, error) {
	//TODO tx rollback yall
	user, err := oa.UserService.GetOrCreateUserByEmail(userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("error AuthUserFlow failed to handle user: %w", err)
	}
	auth, err := oa.AuthService.UpdateOrCreateAuthWithSourceID(user.ID, userInfo.Sub, token, userInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("error AuthUserFlow failed to handle auth: %w", err)
	}
	return auth, user, nil
}
