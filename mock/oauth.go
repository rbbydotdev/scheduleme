package mock

import (
	"net/http"
	"scheduleme/models"
	"scheduleme/oauthhandlers"
	"scheduleme/secure_cookie"

	"golang.org/x/oauth2"
)

type OAuthService struct {
	UserService            models.UserServiceInterface
	EventService           models.EventServiceInterface
	AuthService            models.AuthServiceInterface
	SecureCookie           secure_cookie.SecureCookieInterface
	GoogleAuth             oauthhandlers.GoogleAuthInterface
	HTTPClient             *http.Client
	AuthUserFn             func(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error)
	HandleGoogleCallbackFn func(w http.ResponseWriter, r *http.Request)
	HandleGoogleLoginFn    func(w http.ResponseWriter, r *http.Request)
	HandleLoginSuccessFn   func(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User)
	HandleLogoutFn         func(w http.ResponseWriter, r *http.Request)
}

func (s *OAuthService) AuthUser(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error) {
	return s.AuthUserFn(userInfo, token)
}

func (s *OAuthService) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	s.HandleGoogleCallbackFn(w, r)
}

func (s *OAuthService) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	s.HandleGoogleLoginFn(w, r)
}

func (s *OAuthService) HandleLoginSuccess(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User) {
	s.HandleLoginSuccessFn(w, r, auth, user)
}

func (s *OAuthService) HandleLogout(w http.ResponseWriter, r *http.Request) {
	s.HandleLogoutFn(w, r)
}
