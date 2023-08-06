package oauthhandlers

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// type Source string

func NewGoogleOAuth(redirectURL string, clientID string, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}
}
