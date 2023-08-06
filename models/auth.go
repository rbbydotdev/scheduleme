package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"scheduleme/toerr"
	"scheduleme/values"
	"strings"

	sq "scheduleme/sqlite"
	"time"

	"golang.org/x/oauth2"
)

type AuthService struct {
	db     *sq.Db
	Source values.OAuthSource
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthServiceInterface = (*AuthService)(nil)

type AuthServiceInterface interface {
	CreateAuth(auth *Auth) (values.ID, error)
	// GetAuthByEmail(email string) (*Auth, error)
	GetAuthByUserID(authID values.ID, source values.OAuthSource) (*Auth, error)
	UpdateOrCreateAuthWithSourceID(userID values.ID, sourceID string, token *oauth2.Token, userInfo *UserInfo) (*Auth, error)
	GetAuthBySourceID(sourceID string) (*Auth, error)
	UpdateAuthBySourceID(sourceID string, accessToken Token, refreshToken Token, expiry time.Time) error
	UpdateAuthByID(ID values.ID, accessToken Token, refreshToken Token, expiry time.Time) error
	// LinkAuthIDToUserID(authID values.ID, userID values.ID) error
	// AttachResourceByID(ID values.OAuthSource, ri *RouteInfo) error
}

type UserInfo struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewAuthService(db *sq.Db) *AuthService {
	return &AuthService{db: db}
}

type Auth struct {
	ID        values.ID `json:"id"`
	UserID    values.ID `json:"userID"`
	User      *User     `json:"-"`
	Source    string    `json:"source"`
	SourceID  string    `json:"sourceID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Avatar       string    `json:"picture"`
	AccessToken  Token     `json:"access_token"`
	RefreshToken Token     `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

type AuthToken string

// Loads auth data for user route, to be used for calendar lookup
func (a *AuthService) AuthForUserRoute(ri *RouteInfo, _ context.Context) error {
	auth, err := a.GetAuthByUserID(ri.User.ID, values.OAuthSourceGoogle)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(fmt.Errorf("auth for user id %v not found", ri.User.ID))
		}
		return err
	}
	ri.Auth = *auth
	return nil
}

func (a *AuthService) CreateAuth(auth *Auth) (values.ID, error) {
	// Insert the auth
	res, err := a.db.Exec(`
		INSERT INTO auths (user_id, source_id, source, created_at, updated_at,
			access_token, refresh_token, expiry)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		auth.UserID, auth.SourceID, auth.Source, auth.CreatedAt, auth.UpdatedAt,
		auth.AccessToken, auth.RefreshToken, auth.Expiry)
	if err != nil {
		// Handle the specific error cases
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, toerr.Conflict(fmt.Errorf("auth already exists: %w", err)) //.Msg("auth already exists")
		} else {
			return 0, toerr.Internal(fmt.Errorf("failed to create auth err=%w", err))

		}
	}

	// Get the ID of the new auth
	newAuthID, err := res.LastInsertId()
	if err != nil {
		return 0, toerr.Internal(fmt.Errorf("failed to get new auth ID err=%w", err))
	}

	id := values.ID(newAuthID)
	// Set the ID of the auth
	auth.ID = id
	return id, nil
}

// func (a *AuthService) GetAuthByEmail(email string) (*Auth, error) {
// 	auth := &Auth{}
// 	err := a.db.QueryRow(`
//     		SELECT id, user_id, source_id, source, created_at, updated_at,
//     		       access_token, refresh_token, expiry
//     		FROM auths WHERE email = ?`,
// 		email).Scan(
// 		&auth.ID, &auth.UserID, &auth.SourceID, &auth.Source, &auth.CreatedAt, &auth.UpdatedAt,
// 		&auth.AccessToken, &auth.RefreshToken, &auth.Expiry)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, toerr.NotFound(fmt.Errorf("GetAuthByEmail auth with email %s not found err=%w", email, err)).Msg("auth with email %s not found", email)
// 		}
// 		return nil, toerr.Internal(err)
// 	}
// 	return auth, nil
// }

func (a *AuthService) GetAuthByUserID(userID values.ID, source values.OAuthSource) (*Auth, error) {
	auth := &Auth{}
	//unique because 	, UNIQUE(user_id, source)
	err := a.db.QueryRow(`
    		SELECT id, user_id, source_id, source, created_at, updated_at,
    		       access_token, refresh_token, expiry
    		FROM auths WHERE user_id = ? AND source = ?`,
		userID, source).Scan(
		&auth.ID, &auth.UserID, &auth.SourceID, &auth.Source, &auth.CreatedAt, &auth.UpdatedAt,
		&auth.AccessToken, &auth.RefreshToken, &auth.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(
				fmt.Errorf("GetAuthByUserID auth with UserID %v, not found err=%w", userID, err),
			)
		}
		return nil, toerr.Internal(err)
	}
	return auth, nil
}

func (a *AuthService) UpdateOrCreateAuthWithSourceID(userID values.ID, sourceID string, token *oauth2.Token, userInfo *UserInfo) (*Auth, error) {
	auth, err := a.GetAuthBySourceID(sourceID)
	if err == nil {
		err = a.UpdateAuthBySourceID(auth.SourceID, auth.AccessToken, auth.RefreshToken, auth.Expiry)
	} else if errors.Is(err, sql.ErrNoRows) {
		newAuth := &Auth{
			SourceID:     sourceID,
			Source:       string(a.Source),
			AccessToken:  Token(token.AccessToken),
			Avatar:       userInfo.Picture,
			Name:         userInfo.Name,
			RefreshToken: Token(token.RefreshToken),
			Expiry:       token.Expiry,
			Email:        userInfo.Email,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			UserID:       userID,
		}
		_, err = a.CreateAuth(newAuth)

		auth = newAuth
	}
	return auth, err
}

func (a *AuthService) GetAuthBySourceID(sourceID string) (*Auth, error) {
	auth := &Auth{}
	err := a.db.QueryRow(`
    		SELECT id, user_id, source_id, source, created_at, updated_at,
    		       access_token, refresh_token, expiry
    		FROM auths WHERE source_id = ?`,
		sourceID).Scan(
		&auth.ID, &auth.UserID, &auth.SourceID, &auth.Source, &auth.CreatedAt, &auth.UpdatedAt,
		&auth.AccessToken, &auth.RefreshToken, &auth.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			// Auth not found, return a nil auth and a custom error
			return nil, toerr.NotFound(
				fmt.Errorf("GetAuthBySourceID auth with sourceID %s not found err=%w", sourceID, err),
			) //.Msg("sourceID not found")
		}
		// A different error happened
		return nil, err
	}
	return auth, nil
}

func (a *AuthService) UpdateAuthBySourceID(sourceID string, accessToken Token, refreshToken Token, expiry time.Time) error {
	// Update the auth
	result, err := a.db.Exec(`
    		UPDATE auths SET
    			updated_at = ?,
    			access_token = ?,
    			refresh_token = ?,
    			expiry = ?
    		WHERE source_id = ?`,
		time.Now(), accessToken, refreshToken, expiry, sourceID)
	if err != nil {
		return toerr.Internal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return toerr.Internal(err)
	}
	if rowsAffected == 0 {
		return toerr.NotFound(
			fmt.Errorf("UpdateAuthBySourceID auth with sourceID %s not found err=%w", sourceID, err),
		) //.Msg("sourceID %v not found", sourceID)
	}
	return nil
}

func (a *AuthService) UpdateAuthByID(ID values.ID, accessToken Token, refreshToken Token, expiry time.Time) error {
	result, err := a.db.Exec(`
    		UPDATE auths SET
    			updated_at = ?,
    			access_token = ?,
    			refresh_token = ?,
    			expiry = ?
    		WHERE id = ?`,
		time.Now(), accessToken, refreshToken, expiry, ID)
	if err != nil {
		return toerr.Internal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return toerr.Internal(err)
	}
	if rowsAffected == 0 {
		return toerr.NotFound(
			fmt.Errorf("UpdateAuthByID auth with id=%v not found", ID),
		) //.Msg("auth id %v not found", ID)
	}
	return nil
}
