package toerr_test

import (
	"database/sql"
	"errors"
	"fmt"
	"scheduleme/toerr"
	"testing"
)

func TestWrapAndUnwrapError(t *testing.T) {
	testID := "12345"
	err := sql.ErrNoRows

	wrappedErr := toerr.NotFound(
		fmt.Errorf("testID %s not found err=%w", testID, err),
	).Msg("testID not found")

	// We want to assert that the resulting error.Is(sql.ErrNoRows)
	if !errors.Is(wrappedErr, sql.ErrNoRows) {
		t.Errorf("Expected the wrappedErr to be sql.ErrNoRows but it's not")
	}
}
