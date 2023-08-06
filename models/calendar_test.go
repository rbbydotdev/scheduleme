package models_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"context"
	"scheduleme/hof"
	"scheduleme/mock"
	"scheduleme/models"
	"scheduleme/values"
)

func TestGetBusyTimes(t *testing.T) {

	//This is the response expected from Google Calendar's FreeBusy API.

	client := mock.HTTPClientJSONString(`{
		"kind": "calendar#freeBusy",
		"timeMin": "2022-02-01T00:00:00Z",
		"timeMax": "2022-02-01T23:59:59Z",
		"calendars": {
			"primary": {
				"busy": [
					{
						"start": "2022-02-01T06:00:00Z",
						"end": "2022-02-01T07:00:00Z"
					}
				]
			}
		}
	}`)
	agent := models.NewGoogleAgent(client)

	calendarID := "primary"
	minTime, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
	maxTime, _ := time.Parse(time.RFC3339, "2022-02-01T23:59:59Z")

	slots, err := agent.GetBusyTimes(calendarID, minTime, maxTime)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(slots) != 1 {
		t.Fatalf("Expected one busy slot, but got %d", len(slots))
	}

	expectedStart, _ := time.Parse(time.RFC3339, "2022-02-01T06:00:00Z")
	expectedEnd, _ := time.Parse(time.RFC3339, "2022-02-01T07:00:00Z")

	if !slots[0].Start.Equal(expectedStart) || !slots[0].End.Equal(expectedEnd) {
		t.Fatalf("Expected slot to be from %v to %v, but got from %v to %v",
			expectedStart, expectedEnd, slots[0].Start, slots[0].End)
	}
}

type ReqBody struct {
	CalendarExpansionMax int `json:"calendarExpansionMax"`
	GroupExpansionMax    int `json:"groupExpansionMax"`
	Items                []struct {
		ID string `json:"id"`
	} `json:"items"`
	TimeMin string `json:"timeMin"`
	TimeMax string `json:"timeMax"`
}

func TestEventSlotsForAuth(t *testing.T) {
	moc := mock.OAuth2ClientResponse(`{
		"kind": "calendar#freeBusy",
		"timeMin": "2022-02-01T00:00:00Z",
		"timeMax": "2022-02-01T23:59:59Z",
		"calendars": {
			"primary": {
				"busy": [
					{
						"start": "2022-02-01T06:00:00Z",
						"end": "2022-02-01T07:00:00Z"
					}
				]
			}
		}
	}`)

	gcs := models.GoogleCalendarService{
		OAuth2: moc,
	}

	minTime, _ := time.Parse(time.RFC3339, "2022-03-01T00:00:00Z")
	maxTime, _ := time.Parse(time.RFC3339, "2022-03-30T23:59:59Z")
	auth := &models.Auth{
		RefreshToken: "refresh",
		AccessToken:  "access",
		Expiry:       time.Now().Add(time.Hour),
	}
	slot := values.DurationSlot{
		Start: 12 * time.Hour, // 12 hours from the start of the day : 12pm
		End:   17 * time.Hour, // 17 hours from the start of the day : 4pm
	}
	incm := values.NewIncMask(&values.DurationSlots{slot}, nil, nil)
	event := &models.Event{
		Duration:   2 * time.Hour,
		AvailMasks: &values.AvailMasks{incm},
	}
	dss, err := gcs.EventSlotsForAuth(context.Background(), event, auth, minTime, maxTime)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if hof.Any(*dss, func(ds values.DateSlot) bool {
		return !slot.Withholds(&ds)
	}) {
		t.Fatal(`Expected slots to be withheld in values.DurationSlot{...}`)
	}

	req := moc.GetRequest(0)
	defer req.Body.Close()
	bodyBytes, _ := io.ReadAll(req.Body)

	// Convert the actual request body to a ReqBody object
	var got ReqBody
	err = json.Unmarshal(bodyBytes, &got)
	if err != nil {
		t.Fatalf("couldn't convert request body to JSON: %s", err)
	}

	wantStr := `{
		"calendarExpansionMax": 2,
		"groupExpansionMax": 2,
		"items": [
			{
				"id": "primary"
			}
		],
		"timeMax": "2022-03-30T23:59:59Z",
		"timeMin": "2022-03-01T00:00:00Z"
	}`
	// Convert the 'want' string to a ReqBody object
	var want ReqBody
	err = json.Unmarshal([]byte(wantStr), &want)
	if err != nil {
		t.Fatalf("could not convert want string to JSON: %s", err)
	}

	// Compare got and want
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %+v, but got %+v", want, got)
	}

}
