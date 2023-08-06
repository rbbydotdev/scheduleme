package toerr

/*

Custom app wide errors which map rather well to the logic and flow of the app,
these then could be translated to HTTP status codes and messages.


.Msg is used to add user friendly messages, good especially for ommiting sensitive information

*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"scheduleme/values"
	// "scheduleme/models"
)

func (e *Error) Unwrap() error {
	return e.Err
}

func Render(w http.ResponseWriter, r *http.Request, err error) *Error {
	// Type assertion to check if it's of type *Error
	if e, ok := err.(*Error); ok {
		w.WriteHeader(codes[e.Code])

		// rqi := frame.FromContext[models.RequestInfo](r.Context())
		// if rqi.IsHTML() {
		// 	w.Header().Set("Content-Type", "text/html")
		// 	fmt.Fprint(w, e.ToHTML()) // Writes HTML response
		// } else {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	json.NewEncoder(w).Encode(e) // Marshals error to JSON
		// }
		return e
	}

	// Not an *Error, but some other error type
	w.WriteHeader(http.StatusInternalServerError)
	return Internal(err)
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

//TODO add context / request id to error

type Error struct {
	Code    string
	Message string
	Err     error
	Debug   string
	// Meta    interface{}
	// Op      string
}

var codes = map[string]int{
	ECONFLICT:       http.StatusConflict,
	EINVALID:        http.StatusBadRequest,
	ENOTFOUND:       http.StatusNotFound,
	ENOTIMPLEMENTED: http.StatusNotImplemented,
	EUNAUTHORIZED:   http.StatusUnauthorized,
	EINTERNAL:       http.StatusInternalServerError,
}

const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

// Marshal to JSON without the original error.
func (e *Error) MarshalJSON() ([]byte, error) {
	structErr := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	}
	return json.Marshal(structErr)
}
func (e *Error) ToHTML() string {
	// Escape string for safe HTML representation
	escapedMessage := template.HTMLEscapeString(e.Message)

	html := fmt.Sprintf("<p><strong>Error:</strong> %s</p><p><strong>Message:</strong> %s</p>", e.Code, escapedMessage)
	return html
}

func NewError(code string, err error) *Error {
	if err == nil {
		err = errors.New(code)
	}
	return &Error{
		Code:    code,
		Err:     err,
		Message: code,
	}
}

func Invalid(err error) *Error {
	return NewError(EINVALID, err)
}

// Same Same but diff ;)
func BadRequest(err error) *Error {
	return NewError(EINVALID, err)
}

func Internal(err error) *Error {
	return NewError(EINTERNAL, err)
}

func Conflict(err error) *Error {
	return NewError(ECONFLICT, err)
}

func NotFound(err error) *Error {
	return NewError(ENOTFOUND, err)
}

func IDNotFound(id values.ID, err error) *Error {
	return NewError(ENOTFOUND, fmt.Errorf("id:%v,%w", id, err))
}

func Unauthorized(err error) *Error {
	return NewError(EUNAUTHORIZED, err)
}

func (e *Error) Msg(format string, a ...interface{}) *Error {
	e.Message += ", " + fmt.Sprintf(format, a...)
	return e
}

func (e *Error) Dbg(format string, a ...interface{}) *Error {
	e.Debug += ", " + fmt.Sprintf(format, a...)
	return e
}

type Config struct {
	logger *log.Logger
	// omitMsgValues bool
}

var config = Config{
	logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	// omitMsgValues: true,
}

func Init(cfg *Config) {
	if cfg != nil {
		config = *cfg
	}
}

func (e *Error) Error() string {
	// Do not include e.Err in the error string
	// which should be kept internal or logged at server side.
	return fmt.Sprintf("code= %s message= %s err= %v debug=%v", e.Code, e.Message, e.Err, e.Debug)
}
func (e *Error) Log() *Error {
	// config.logger.Printf("code= %s message= %s err= %v debug= %v", e.Code, e.Message, e.Err, e.Debug)
	config.logger.Print(e.Error())
	return e
}

func (e *Error) Render(w http.ResponseWriter, r *http.Request) *Error {
	Render(w, r, e)
	return e
}
