package mock

import (
	"net/http"
	"scheduleme/models"
	"scheduleme/secure_cookie"
	"scheduleme/values"
)

var _ secure_cookie.SecureCookieInterface = (*SecureCookie)(nil)

type SecureCookie struct {
	ClearSessionFn       func(w http.ResponseWriter)
	GetOrCreateSessionFn func(r *http.Request) *models.SessionInfo
	GetSessionFn         func(r *http.Request) (*models.SessionInfo, error)
	HttpRedirectFn       func(w http.ResponseWriter, r *http.Request, url string, status int)
	NextServeHTTPFn      func(next http.Handler, w http.ResponseWriter, r *http.Request)
	PopRedirectURLFn     func(r *http.Request) string
	PushRedirectURLFn    func(r *http.Request, url string)
	SetSessionFn         func(w http.ResponseWriter, r *http.Request, session *models.SessionInfo) error
	SetSessionFlashFn    func(w http.ResponseWriter, r *http.Request, flash string) error
	SetSessionStateFn    func(w http.ResponseWriter, r *http.Request, state string) error
	SetSessionUserIDFn   func(w http.ResponseWriter, r *http.Request, userID values.ID) error
	DeserializeSessionFn func(sessionStr string) (*models.SessionInfo, error)
	SerializeSessionFn   func(session *models.SessionInfo) (string, error)
}

func (sc *SecureCookie) SerializeSession(session *models.SessionInfo) (string, error) {
	return sc.SerializeSessionFn(session)
}
func (sc *SecureCookie) DeserializeSession(sessionStr string) (*models.SessionInfo, error) {
	return sc.DeserializeSessionFn(sessionStr)
}
func (sc *SecureCookie) SetSession(w http.ResponseWriter, r *http.Request, session *models.SessionInfo) error {
	return sc.SetSessionFn(w, r, session)
}
func (sc *SecureCookie) GetSession(r *http.Request) (*models.SessionInfo, error) {
	return sc.GetSessionFn(r)
}
func (sc *SecureCookie) ClearSession(w http.ResponseWriter) {
	sc.ClearSessionFn(w)
}
func (sc *SecureCookie) GetOrCreateSession(r *http.Request) *models.SessionInfo {
	return sc.GetOrCreateSessionFn(r)
}

func (sc *SecureCookie) HttpRedirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	sc.HttpRedirectFn(w, r, url, status)
}

func (sc *SecureCookie) NextServeHTTP(next http.Handler, w http.ResponseWriter, r *http.Request) {
	sc.NextServeHTTPFn(next, w, r)
}

func (sc *SecureCookie) PopRedirectURL(r *http.Request) string {
	return sc.PopRedirectURLFn(r)
}
func (sc *SecureCookie) PushRedirectURL(r *http.Request, url string) {
	sc.PushRedirectURLFn(r, url)
}

func (sc *SecureCookie) SetSessionFlash(w http.ResponseWriter, r *http.Request, flash string) error {
	return sc.SetSessionFlashFn(w, r, flash)
}

func (sc *SecureCookie) SetSessionState(w http.ResponseWriter, r *http.Request, state string) error {
	return sc.SetSessionStateFn(w, r, state)
}
func (sc *SecureCookie) SetSessionUserID(w http.ResponseWriter, r *http.Request, userID values.ID) error {
	return sc.SetSessionUserIDFn(w, r, userID)
}
