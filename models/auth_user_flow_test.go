package models_test

import (
	"scheduleme/models"
	"scheduleme/sqlite"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// func _TestAuthUserFlow(t *testing.T) {

// 	as := &mock.AuthService{}
// 	us := &mock.UserService{}

// 	auf := models.NewAuthUserFlow(as, us)

// 	ui := &models.UserInfo{}
// 	tok := &models.Token{}

// 	auf.AuthUser(nil, nil)

// }

func TestAuthUserFlow(t *testing.T) {

	db, err := sqlite.NewOpenDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	as := models.NewAuthService(db)
	us := models.NewUserService(db)

	auf := models.NewAuthUserFlow(as, us)

	ui := &models.UserInfo{
		Sub:     "123",
		Name:    "test",
		Email:   "foobar@bizzbazz",
		Picture: "https://foobar.com",
	}
	tok := &oauth2.Token{
		AccessToken:  "123",
		RefreshToken: "123",
		Expiry:       time.Now().Add(time.Hour),
	}

	_, _, err = auf.AuthUser(ui, tok)
	if err != nil {
		t.Fatalf("Failed to auth user: %v", err)
	}

}
