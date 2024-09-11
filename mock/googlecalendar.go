package mock

import (
	"context"
	"io"
	"log"
	"net/http"
	"scheduleme/models"
	"scheduleme/values"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type GoogleCalendarService struct {
	CalendarClientFn    func(ctx context.Context, a *models.Auth) *http.Client
	EventSlotsForAuthFn func(ctx context.Context, e *models.Event, a *models.Auth, startTime time.Time, endTime time.Time) (*values.DateSlots, error)
	BuildAgentFn        func(context.Context, *models.Auth) models.OAuthAgent
}

func (s *GoogleCalendarService) CalendarClient(ctx context.Context, a *models.Auth) *http.Client {
	return s.CalendarClientFn(ctx, a)
}
func (s *GoogleCalendarService) EventSlotsForAuth(ctx context.Context, e *models.Event, a *models.Auth, startTime time.Time, endTime time.Time) (*values.DateSlots, error) {
	return s.EventSlotsForAuthFn(ctx, e, a, startTime, endTime)
}
func (s *GoogleCalendarService) BuildAgent(ctx context.Context, a *models.Auth) models.OAuthAgent {
	return s.BuildAgentFn(ctx, a)
}

/*
Usage:

	moc := &mock.MockOAuth2Clientable{Response: `{
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
	}`}

	gcs := models.GoogleCalendarService{
		OAuth2: moc,
	}

gcs.EventSlotsForAuth(context.Background(), &models.Event{}, &models.Auth{}, time.Now(), time.Now())
*/
type MockOAuth2Clientable struct {
	Response string
}

func (moc *MockOAuth2Clientable) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return HTTPClientJSONString(moc.Response)
}

type ClientFunc struct {
	f              func(ctx context.Context, t *oauth2.Token) *http.Client
	ReqStore       []*http.Request
	CustomResponse string
}

func (cf *ClientFunc) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	client := cf.f(ctx, t)
	client.Transport = cf
	return client
}

func (cf *ClientFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	cf.ReqStore = append(cf.ReqStore, req)
	// Use Custom Response stored in ClientFunc while returning Response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(cf.CustomResponse)),
	}, nil
}

// Add this function to retrieve a Request
func (cf *ClientFunc) GetRequest(index int) *http.Request {
	if index < len(cf.ReqStore) {
		return cf.ReqStore[index]
	}
	log.Println("Request not found for index:", index)
	return nil
}

func OAuth2ClientResponse(res string) *ClientFunc {
	return &ClientFunc{
		ReqStore:       make([]*http.Request, 0),
		CustomResponse: res,
		f: func(ctx context.Context, t *oauth2.Token) *http.Client {
			return HTTPClientJSONString(res)
		},
	}
}
