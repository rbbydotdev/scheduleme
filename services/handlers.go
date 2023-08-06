package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scheduleme/frame"
	"scheduleme/models"
	"scheduleme/secure_cookie"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func HandleSessionContextDebug(w http.ResponseWriter, r *http.Request) {
	ses := frame.FromContext[models.SessionInfo](r.Context())
	bytes, err := json.MarshalIndent(ses, "", "\t")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Fprintf(w, "%s\n", bytes)
}

func HandleContextDebug(w http.ResponseWriter, r *http.Request) {
	ri := frame.FromContext[models.RouteInfo](r.Context())
	bytes, err := json.MarshalIndent(ri, "", "\t")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Fprintf(w, "%s\n", bytes)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {

	ses := frame.FromContext[models.SessionInfo](r.Context())
	if ses.UserID == 0 {
		fmt.Fprintln(w, `<a href="/login">Login with Google</a>`)
	} else {
		fmt.Fprintf(w, `
    <h1>Hello User %v</h1>
    <a href="/logout?csrfToken=%s">Logout</a>
    `, ses.UserID, ses.CSRFToken)
	}
}
func RedirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusFound)
}

func HandleSessionDebug(sc *secure_cookie.SecureCookie) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ses, _ := sc.GetSession(r)
		bytes, err := json.MarshalIndent(ses, "", "\t")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Fprintf(w, "%s\n", bytes)
	}
}
