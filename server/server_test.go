package server_test

import (
	"scheduleme/server"
	"testing"
)

// TestRunMain tests the RunMain function

func TestRunMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RunMain panicked with: %v", r)
		}
	}()
	server.RunMain()
}
