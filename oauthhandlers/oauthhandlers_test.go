package oauthhandlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"scheduleme/frame"
	"scheduleme/mock"
	"scheduleme/models"

	"scheduleme/oauthhandlers"

	// "scheduleme/secure_cookie"
	"testing"

	"golang.org/x/oauth2"
)

func TestHandleGoogleCallback(t *testing.T) {
	mockClient := mock.HTTPClientJSONBody(map[string]interface{}{
		"sub":   "dummySub",
		"name":  "dummyName",
		"email": "dummyEmail",
	})
	authUserFlow := &mock.AuthUserFlow{}
	googleAuth := mock.NewGoogleAuth("/?code=dummyCode&state=dummyState", "dummyToken")
	googleAuth.ExchangeFn = func(ctx context.Context, code string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
		return &oauth2.Token{}, nil
	}
	secureCookie := &mock.SecureCookie{}

	authUserFlow.AuthUserFn = func(userInfo *models.UserInfo, token *oauth2.Token) (*models.Auth, *models.User, error) {
		user := &models.User{}
		auth := &models.Auth{}
		return auth, user, nil
	}

	oauth := oauthhandlers.NewOAuth(
		secureCookie,
		googleAuth,
		authUserFlow,
		mockClient,
	)

	recreq := func() (*httptest.ResponseRecorder, *http.Request) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/?code=dummyCode&state=dummyState", nil)
		if err != nil {
			t.Fatal(err)
		}
		return rec, req
	}

	recorder, req := recreq()
	sessionCtx := frame.NewContextWith(req.Context(), &models.SessionInfo{State: "______"})
	oauth.HandleGoogleCallback(
		func(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User) {
		})(recorder, req.WithContext(sessionCtx))
	if got, want := recorder.Code, http.StatusBadRequest; got != want {
		t.Errorf("State MISMATCH should give http.StatusBadRequest")
	}

	recorder, req = recreq()
	sessionCtx = frame.NewContextWith(req.Context(), &models.SessionInfo{State: "dummyState"})
	oauth.HandleGoogleCallback(
		func(w http.ResponseWriter, r *http.Request, auth *models.Auth, user *models.User) {
		})(recorder, req.WithContext(sessionCtx))
	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Errorf("Expecting http.StatusAccepted, got %v want %v", got, want)
	}

}
