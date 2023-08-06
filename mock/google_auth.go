package mock

import (
	"context"

	"golang.org/x/oauth2"
)

type GoogleAuth struct {
	authCodeURL   string
	accessToken   string
	AuthCodeURLFn func(string, ...oauth2.AuthCodeOption) string
	ExchangeFn    func(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

func NewGoogleAuth(authCodeURL string, accessToken string) *GoogleAuth {
	return &GoogleAuth{
		authCodeURL: authCodeURL,
		accessToken: accessToken,
	}
}

func (m *GoogleAuth) Exchange(ctx context.Context, code string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return m.ExchangeFn(ctx, code)
}

func (m *GoogleAuth) AuthCodeURL(nonce string, _ ...oauth2.AuthCodeOption) string {
	return m.AuthCodeURLFn(nonce)
}
