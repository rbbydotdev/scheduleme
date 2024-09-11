package server_test

import (
	"scheduleme/server"
	"testing"
)

// TestRunMain tests the RunMain function

func TestRun(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run panicked with: %v", r)
		}
	}()
	server.Run(nil, nil)
}
