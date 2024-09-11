package dbutil

import "database/sql"

func WithCount(result sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	return count, err
}
