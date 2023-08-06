package models

import (
	"database/sql"
	"scheduleme/values"
)

type Token = values.Token

// Helper function to returns the number of rows affected by the query.
func withCount(result sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	return count, err
}
