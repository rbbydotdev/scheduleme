package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"scheduleme/models"
	"scheduleme/server"
	"scheduleme/sqlite"
	"scheduleme/test"
	"testing"
)

// 	return &ConfigStruct{
// 		GoogleClientSecret: "mock-secret",
// 		GoogleClientID:     "mock-id",
// 		GoogleRedirectURL:  "http://mockurl",
// 		Port:               "8080",
// 		ENV:                "test",
// 		Dsn:                "mock-dsn",
// 		Secret:             "mock-secret",
// 	}
// }

func TestIntegration(t *testing.T) {
	db, err := sqlite.NewOpenDB(":memory:")
	if err != nil {
		return
	}
	// Run your integration tests here
	err = test.Seed(db)
	if err != nil {
		t.Fatal(err)
	}
	tstCfg := test.TestConfig()
	tstCfg.Port = "0"

	addrStr := server.Run(tstCfg, db)

	t.Run("Get Home", func(t *testing.T) {
		resp, err := http.Get("http://" + addrStr) // Use the server address for the HTTP call
		//check status code
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if err != nil {
			t.Fatalf("Failed to make HTTP GET request: %v", err)
		}
		defer resp.Body.Close()
	})
	t.Run("Get User 1", func(t *testing.T) {
		resp, err := http.Get("http://" + addrStr + "/api/users/1") // Use the server address for the HTTP call
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		if err != nil {
			t.Fatalf("Failed to make HTTP GET request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		var user models.UserView
		err = json.Unmarshal(body, &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON response: %v", err)
		}

		// Expected user
		expectedUser := models.UserView{ID: 1, Name: "test1"}

		if user != expectedUser {
			t.Fatalf("Expected user to be %+v, got %+v", expectedUser, user)
		}
	})

}
