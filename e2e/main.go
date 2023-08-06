package e2e_test

import (
	"testing"
)

// func TestMain(m *testing.M) {
// 	// Do setup here
// 	m.Run()
// 	// Do teardown here
// }

func TestE2E(t *testing.T) {
	// Run your end-to-end tests here
	t.Run("Test1", func(t *testing.T) {
		t.Log("Test1")
	})
	t.Run("Test2", func(t *testing.T) {
		t.Log("Test2")
	})

}
