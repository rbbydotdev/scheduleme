package secure_cookie_test

import (
	// "fmt"

	"net/http"
	"net/http/httptest"
	"scheduleme/models"
	sc "scheduleme/secure_cookie"
	"scheduleme/values"
	"strings"
	"testing"
)

func TestSerializeAndDeserializeSession(t *testing.T) {
	// Arrange
	secret := "abc123"
	cookieName := "testCookie"
	secureCookie := sc.NewSecureCookie(secret, cookieName)

	session := &models.SessionInfo{
		UserID: values.ID(1),
	}

	// Act
	serialized, err := secureCookie.SerializeSession(session)
	if err != nil {
		t.Errorf("Failed to serialize the session: %v", err)
		return
	}
	// Now test deserialization
	deserialized, err := secureCookie.DeserializeSession(serialized)
	if err != nil {
		t.Errorf("Failed to deserialize the session: %v", err)
		return
	}

	// Assert
	// You can assert other properties of the session as per the requirements...
	if deserialized.UserID != session.UserID {
		t.Errorf("The UserID of the deserialized session does not match the original UserID")
	}
}

func TestDeserializeSessionFail(t *testing.T) {
	// Arrange

	evilSecureCookie := sc.NewSecureCookie("evil", "testCookie")
	goodSecureCookie := sc.NewSecureCookie("good", "testCookie")

	session := &models.SessionInfo{
		UserID: values.ID(1),
	}

	// Act
	evilSerialized, err := evilSecureCookie.SerializeSession(session)
	if err != nil {
		t.Errorf("Failed to serialize the session: %v", err)
		return
	}
	// Now test deserialization
	_, err = goodSecureCookie.DeserializeSession(evilSerialized)
	if err == nil {
		t.Errorf("Failed to return err on evil session: %v", err)
		return
	}
}

func TestExtractVerifySession(t *testing.T) {
	// Create a new secure cookie
	secret := "secretpassword"
	cookieName := "testCookie"
	secureCookie := sc.NewSecureCookie(secret, cookieName)

	// Create a new session
	session := &models.SessionInfo{
		UserID: values.ID(1),
		// Add other session properties here...
	}

	// Serialize the session
	sessionStr, err := secureCookie.SerializeSession(session)
	if err != nil {
		t.Fatal(err)
	}

	// Call ExtractVerifySession
	ses, err := secureCookie.ExtractVerifySession(sessionStr)

	// The returned session should not be nil as we passed in a valid sessionStr
	if ses == nil {
		t.Errorf("Expected a session, but got nil")
	}

	// The returned error should be nil as we passed in a valid sessionStr
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	// Test some invalid inputs
	_, err = secureCookie.ExtractVerifySession("")
	if err == nil {
		t.Errorf("Expected an error for empty string, but got nil")
	}

	_, err = secureCookie.ExtractVerifySession("invalidData")
	if err == nil {
		t.Errorf("Expected an error for invalid string, but got nil")
	}

	// Test an 'evil' session with a forged signature.
	evilCookie := sc.NewSecureCookie("wrongSecret", cookieName)
	evilSessionStr, err := evilCookie.SerializeSession(session)
	if err != nil {
		t.Fatal(err)
	}

	_, err = secureCookie.ExtractVerifySession(evilSessionStr)
	if err == nil || !strings.Contains(err.Error(), "invalid signature") {
		t.Errorf("Expected an error with 'invalid signature', but got %v", err)
	}

}

// Test get existing session
func TestGetSession(t *testing.T) {
	// Create a new secure cookie
	secret := "secretpassword"
	cookieName := "testCookie"
	secureCookie := sc.NewSecureCookie(secret, cookieName)

	// Create a new request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new session
	session := &models.SessionInfo{
		UserID: values.ID(1),
		// add other session properties here...
	}

	// Record the response (this implements http.ResponseWriter)
	w := httptest.NewRecorder()

	// Set the session in the cookie
	err = secureCookie.SetSession(w, r, session)
	if err != nil {
		t.Fatal(err)
	}

	// Extract the cookies from the response
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("Expected at least one cookie, but got none")
		return
	}

	// Find the cookie we're interested in
	var foundCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			foundCookie = cookie
			break
		}
	}
	if foundCookie == nil {
		t.Fatalf("Cookie %s not found in response", cookieName)
	}

	// Setup a new request and add the found cookie to the request
	r, _ = http.NewRequest("GET", "/", nil)
	r.AddCookie(foundCookie)

	// Get the session and test it
	gotSession, err := secureCookie.GetSession(r)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if gotSession.UserID != session.UserID {
		t.Errorf("Got session with ID: %v; Expected: %v", gotSession.UserID, session.UserID)
	}
}

// Test ClearSession function
func TestClearSession(t *testing.T) {
	secret := "secretpassword"
	cookieName := "testCookie"
	secureCookie := sc.NewSecureCookie(secret, cookieName)

	// Assert that ClearSession doesn't return an error
	w := httptest.NewRecorder() // An http.ResponseWriter for testing purposes
	secureCookie.ClearSession(w)
	result := w.Result()
	// We're expecting result.StatusCode to be http.StatusOK (i.e., 200)
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected http.StatusOK but received %d", result.StatusCode)
	}
}
