package values

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	// "fmt"
)

type ID int64

type APIKey string

type Token string

type CtxState string

type OAuthSource string

const OAuthSourceGoogle OAuthSource = "google"

func (ts *DateSlots) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

func (ts *DateSlots) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i []DateSlot
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*ts = i
	return nil
}

func (c CtxState) CompareStates(state string) bool {
	return c != "" && c == CtxState(state)
}

func (a *APIKey) String() string {
	return "[REDACTED]"
}

func (t *Token) String() string {
	return "[REDACTED]"
}

func (a APIKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
