package oauthhandlers

import (
	"fmt"
	"scheduleme/secure_cookie"
	"scheduleme/toerr"

	"context"
	"encoding/json"
	"net/http"
	"scheduleme/models"
	"scheduleme/util"

	"scheduleme/frame"

	"golang.org/x/oauth2"
)

type OAuth struct {
	AuthService  models.AuthServiceInterface
	SecureCookie secure_cookie.SecureCookieInterface
	AuthUserFlow models.AuthUserFlowInterface
	GoogleAuth   GoogleAuthInterface
	HTTPClient   *http.Client
}

type GoogleAuthInterface interface {
	Exchange(ctx context.Context, code string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
}

type OAuthInterface interface {
	AuthUser(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error)
	HandleGoogleCallback(w http.ResponseWriter, r *http.Request)
	HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
	HandleLoginSuccess(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User)
	HandleLogout(w http.ResponseWriter, r *http.Request)
}

func NewOAuth(
	sc secure_cookie.SecureCookieInterface,
	ga GoogleAuthInterface,
	auf models.AuthUserFlowInterface,
	ht *http.Client,
) *OAuth {
	return &OAuth{
		SecureCookie: sc,
		GoogleAuth:   ga,
		AuthUserFlow: auf,
		HTTPClient:   ht,
	}
}

func (oa *OAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	queryCSRF := r.URL.Query().Get("csrfToken")
	frame.ModifyContextWith(r.Context(), func(ses *models.SessionInfo) {
		if ses.CSRFToken == queryCSRF {
			*ses = *models.NewSession()
		}
	})

	oa.SecureCookie.HttpRedirect(w, r, "/home", http.StatusFound)
}

func (oa *OAuth) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	nonce := util.RandomStr(32)

	frame.ModifyContextWith(r.Context(), func(ses *models.SessionInfo) {
		ses.State = nonce
	})
	url := oa.GoogleAuth.AuthCodeURL(nonce, oauth2.AccessTypeOffline)
	oa.SecureCookie.HttpRedirect(w, r, url, http.StatusTemporaryRedirect)
}

func (oa *OAuth) HandleLoginSuccess(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User) {
	redir := oa.SecureCookie.PopRedirectURL(r)
	if redir == "" {
		redir = "/home"
	}
	frame.ModifyContextWith(r.Context(), func(ses *models.SessionInfo) {
		ses.UserID = user.ID
		ses.IsAdmin = user.IsAdmin
	})
	oa.SecureCookie.HttpRedirect(w, r, redir, http.StatusFound)
}

func (oa *OAuth) HandleGoogleCallback(handleLoginSuccessFn func(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code, state := r.URL.Query().Get("code"), r.URL.Query().Get("state")
		// ctxState := frame.CtxState(frame.FromContext[models.SessionInfo](r.Context()).State)
		ctxState := models.StateFromContext(r.Context())
		if !ctxState.CompareStates(state) {
			toerr.BadRequest(
				fmt.Errorf("HandleGoogleCallback state mismatch"),
			).Msg("state mismatch").Render(w, r)
			return
		}

		token, err := oa.GoogleAuth.Exchange(context.Background(), code)
		if err != nil {
			toerr.Internal(
				fmt.Errorf("HandleGoogleCallback failed to exchange token %w", err),
			).Msg(
				"failed to exchange token",
			).Render(w, r)
			return
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			toerr.Internal(
				fmt.Errorf("HandleGoogleCallback failed to get user info: %w", err),
			).Msg(
				"failed to get user info",
			).Render(w, r)
			return
		}
		defer response.Body.Close()

		var userInfo models.UserInfo
		if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
			toerr.Internal(
				fmt.Errorf("HandleGoogleCallback failed to decode user info: %w", err),
			).Msg(
				"failed to decode user info",
			).Render(w, r)
			return
		}

		auth, user, err := oa.AuthUserFlow.AuthUser(&userInfo, token)
		if err != nil {
			toerr.Internal(
				fmt.Errorf("HandleGoogleCallback failed auth user flow: %w", err),
			).Msg(
				"failed auth user flow",
			).Render(w, r)
			return
		}
		handleLoginSuccessFn(w, r, auth, user)
	}
}
