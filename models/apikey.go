package models

import (
	// "database/sql"
	"database/sql"
	sq "scheduleme/sqlite"
	"scheduleme/toerr"
	"scheduleme/util"
	"scheduleme/values"
)

type APIKey struct {
	Key    values.APIKey
	ID     values.ID
	UserID values.ID
}
type APIKeyService struct {
	DB *sq.Db
}

func (a *APIKeyService) GetAPIKeyByID(apiKeyID values.ID, userID values.ID) (*APIKey, error) {
	row := a.DB.QueryRow(`SELECT key FROM api_keys WHERE id = ? AND user_id = ?`, apiKeyID, userID)
	var apiKey APIKey
	err := row.Scan(&apiKey)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

type APIKeyInterface interface {
	GetUserIDByAPIKey(key string) (values.ID, error)
	CreateAPIKeyForUserID(userID values.ID) (*APIKey, error)
	DeleteAPIKeyByID(apiKeyID values.ID) (int64, error)
}

func (a *APIKeyService) AttachRemoteByID(apiKeyID values.ID, ri *RouteInfo) error {
	apiKey, err := a.GetAPIKeyByID(apiKeyID, ri.User.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(err)
		}
		return toerr.Internal(err)
	}
	ri.APIKey = *apiKey
	return nil
}

func NewAPIKeyService(db *sq.Db) *APIKeyService {
	return &APIKeyService{DB: db}
}

func (a *APIKeyService) GetUserIDByAPIKey(key string) (values.ID, error) {
	row := a.DB.QueryRow(`SELECT user_id FROM api_keys WHERE key = ?`, key)
	var userID values.ID
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *APIKeyService) DeleteAPIKeyByID(apiKeyID values.ID) (int64, error) {
	return withCount(a.DB.Exec(`DELETE FROM api_keys WHERE id = ?`, apiKeyID))
}

func (a *APIKeyService) CreateAPIKeyForUserID(userID values.ID) (*APIKey, error) {
	keyStr := generateKey()
	var apiKey APIKey
	res, err := a.DB.Exec(`INSERT INTO api_keys (user_id, key) VALUES (?, ?)`, userID, keyStr)
	if err != nil {
		return nil, err
	}
	resID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	apiKey.ID = values.ID(resID)
	apiKey.Key = values.APIKey(keyStr)
	apiKey.UserID = userID
	return &apiKey, nil
}

func generateKey() string {
	return util.RandomStr(32)
}
