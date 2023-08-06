package secure_cookie

import (
	"bytes"
	"fmt"
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/values"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
)

type ID = values.ID

type SecureCookieInterface interface {
	ClearSession(w http.ResponseWriter)
	GetOrCreateSession(r *http.Request) *models.SessionInfo
	GetSession(r *http.Request) (*models.SessionInfo, error)
	HttpRedirect(w http.ResponseWriter, r *http.Request, url string, status int)
	NextServeHTTP(next http.Handler, w http.ResponseWriter, r *http.Request)
	PopRedirectURL(r *http.Request) string
	PushRedirectURL(r *http.Request, url string)
	SetSession(w http.ResponseWriter, r *http.Request, session *models.SessionInfo) error
	SetSessionFlash(w http.ResponseWriter, r *http.Request, flash string) error
	SetSessionState(w http.ResponseWriter, r *http.Request, state string) error
	SetSessionUserID(w http.ResponseWriter, r *http.Request, userID ID) error
	DeserializeSession(sessionStr string) (*models.SessionInfo, error)
	SerializeSession(session *models.SessionInfo) (string, error)
}

// encodeSessionToGob encodes a session to a gob
func encodeSessionToGob(session *models.SessionInfo) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(session)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// createSignature creates a HMAC signature for a given byte array
func (sc *SecureCookie) createSignature(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(sc.Secret))
	h.Write(data)
	return h.Sum(nil), nil
}

type SecureCookie struct {
	Secret     string
	CookieName string
}

func NewSecureCookie(secret string, cookieName string) *SecureCookie {
	return &SecureCookie{secret, cookieName}
}

// SerializeSession serializes a session into a string
func (sc *SecureCookie) SerializeSession(session *models.SessionInfo) (string, error) {
	ses, err := encodeSessionToGob(session)
	if err != nil {
		return "", err
	}
	sig, err := sc.createSignature(ses)
	if err != nil {
		return "", err
	}
	encSig := base64.StdEncoding.EncodeToString(sig)
	encSes := base64.StdEncoding.EncodeToString(ses)
	return fmt.Sprintf("%s:%s", encSig, encSes), nil
}

func (sc *SecureCookie) ExtractVerifySession(sessionStr string) ([]byte, error) {
	parts := strings.SplitN(sessionStr, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid session string")
	}
	sig, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	ses, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	expectedSig, err := sc.createSignature(ses)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(sig, expectedSig) {
		return nil, fmt.Errorf("invalid signature")
	}
	return ses, nil
}

// DeserializeSession deserializes a session from a string
func (sc *SecureCookie) DeserializeSession(sessionStr string) (*models.SessionInfo, error) {
	ses, err := sc.ExtractVerifySession(sessionStr)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(ses)
	dec := gob.NewDecoder(buf)
	var session models.SessionInfo
	err = dec.Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (sc *SecureCookie) SaveSessionInfoToCookie(w http.ResponseWriter, r *http.Request) {
	ses := frame.FromContext[models.SessionInfo](r.Context())
	sc.SetSession(w, r, ses)
}
func (sc *SecureCookie) SetSession(w http.ResponseWriter, r *http.Request, session *models.SessionInfo) error {
	// frame.NewContextWith(r.Context(), session)
	fullSessionStr, err := sc.SerializeSession(session)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sc.CookieName,
		Value:    fullSessionStr,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour), //sliding window of expiration
	})
	return nil
}

func (sc *SecureCookie) PopRedirectURL(r *http.Request) string {
	sesCtx := frame.FromContext[models.SessionInfo](r.Context())
	redir := sesCtx.RedirectURL
	frame.ModifyContextWith(r.Context(), func(ses *models.SessionInfo) {
		ses.RedirectURL = ""
	})
	return redir
}

func (sc *SecureCookie) PushRedirectURL(r *http.Request, url string) {
	frame.ModifyContextWith(r.Context(), func(ses *models.SessionInfo) {
		ses.RedirectURL = url
	})
}

func (sc *SecureCookie) SetSessionState(w http.ResponseWriter, r *http.Request, state string) error {
	ses := sc.GetOrCreateSession(r)
	ses.State = state
	return sc.SetSession(w, r, ses)
}

func (sc *SecureCookie) SetSessionFlash(w http.ResponseWriter, r *http.Request, flash string) error {
	ses := sc.GetOrCreateSession(r)
	ses.Flash = flash
	return sc.SetSession(w, r, ses)
}

func (sc *SecureCookie) SetSessionUserID(w http.ResponseWriter, r *http.Request, userID ID) error {
	ses := sc.GetOrCreateSession(r)
	ses.UserID = userID
	return sc.SetSession(w, r, ses)
}

func (sc *SecureCookie) ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sc.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})
}

func (sc *SecureCookie) GetSession(r *http.Request) (*models.SessionInfo, error) {
	cookie, err := r.Cookie(sc.CookieName)
	//ErrNoCookie is returned if a cookie is not found.
	if err != nil {
		return nil, err
	}
	return sc.DeserializeSession(cookie.Value)
}

func (sc *SecureCookie) GetOrCreateSession(r *http.Request) *models.SessionInfo {
	//could pull from db here
	ses, err := sc.GetSession(r)
	if err != nil {
		return models.NewSession()
	}
	return ses
}

// Helper to ensure call to first save session context into cookie before continueing to next
func (sc *SecureCookie) NextServeHTTP(next http.Handler, w http.ResponseWriter, r *http.Request) {
	//could write to db here
	sc.SaveSessionInfoToCookie(w, r)
	next.ServeHTTP(w, r)
}

// Helper to first save session context into cookie then Calls http.Redirect
func (sc *SecureCookie) HttpRedirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	sc.SaveSessionInfoToCookie(w, r)
	http.Redirect(w, r, url, status)
}
