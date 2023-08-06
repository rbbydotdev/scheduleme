package models

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"scheduleme/frame"
	"scheduleme/toerr"
	"scheduleme/values"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleAgent struct {
	Client *http.Client
}

func NewGoogleAgent(
	client *http.Client,
) *GoogleAgent {
	return &GoogleAgent{
		Client: client,
	}
}

type OAuth2Clientable interface {
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

type GoogleCalendarService struct {
	OAuth2 OAuth2Clientable
}

type CalendarService interface {
	CalendarClient(ctx context.Context, a *Auth) *http.Client
	BuildAgent(context.Context, *Auth) OAuthAgent
	EventSlotsForAuth(ctx context.Context, e *Event, a *Auth, startTime time.Time, endTime time.Time) (*values.DateSlots, error)
}

func NewGoogleCalendarService(
	oac *oauth2.Config,
) *GoogleCalendarService {
	return &GoogleCalendarService{
		OAuth2: oac,
	}
}

type OAuthAgent interface {
	GetBusyTimes(calendarID string, minTime time.Time, maxTime time.Time) (values.DateSlots, error)
	GetDefaultBusyTimes(minTime time.Time, maxTime time.Time) (values.DateSlots, error)
}

func (cs *GoogleCalendarService) AvailabilityForEventRoute(ri *RouteInfo, ctx context.Context) (err error) {
	qi := frame.FromContext[QueryInfo](ctx)
	slots, err := cs.EventSlotsForAuth(
		ctx,
		&ri.Event,
		&ri.Auth,
		qi.AvailQuery.StartTime,
		qi.AvailQuery.EndTime,
	)
	ri.Availability = *slots
	return
}

func (cs *GoogleCalendarService) BuildAgent(ctx context.Context, auth *Auth) OAuthAgent {
	return &GoogleAgent{cs.CalendarClient(ctx, auth)}
}

func (cs *GoogleCalendarService) EventSlotsForAuth(ctx context.Context, event *Event, auth *Auth, startTime time.Time, endTime time.Time) (*values.DateSlots, error) {
	agent := cs.BuildAgent(ctx, auth)
	ds, err := agent.GetDefaultBusyTimes(startTime, endTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(fmt.Errorf("error EventSlotsForAuth, GetDefaultBusyTimes: %w", err)).Msg("no default busy times found")
		}
		return nil, toerr.Internal(fmt.Errorf("error EventSlotsForAuth, GetDefaultBusyTimes: %w", err)).Msg("error getting default busy times")
	}

	masks := *event.AvailMasks
	busyMask := values.BusyTimesToMask(&ds)
	masks = append(masks, busyMask)
	return masks.GetDateSlots(event.Duration, startTime, endTime), nil
}

func (cs *GoogleCalendarService) CalendarClient(
	ctx context.Context,
	a *Auth,
) *http.Client {
	token := &oauth2.Token{
		AccessToken:  string(a.AccessToken),
		RefreshToken: string(a.RefreshToken),
		Expiry:       a.Expiry,
	}
	return cs.OAuth2.Client(ctx, token)
}

func (g *GoogleAgent) GetDefaultBusyTimes(minTime time.Time, maxTime time.Time) (values.DateSlots, error) {
	return g.GetBusyTimes("primary", minTime, maxTime)
}

func (g *GoogleAgent) GetBusyTimes(calendarID string, minTime time.Time, maxTime time.Time) (values.DateSlots, error) {
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(g.Client))
	if err != nil {
		return nil, toerr.Internal(fmt.Errorf("error GetBusyTimes, unable to create calendar service: %w", err)).Msg("error creating calendar service")
	}

	//TODO: handle and mitigate HTTP error code 429 rate limit reached -
	//https://gist.github.com/MelchiSalins/27c11566184116ec1629a0726e0f9af5
	freeBusyResp, err := srv.Freebusy.Query(&calendar.FreeBusyRequest{
		CalendarExpansionMax: 2,
		GroupExpansionMax:    2,
		TimeMin:              minTime.Format(time.RFC3339),
		TimeMax:              maxTime.Format(time.RFC3339),
		Items: []*calendar.FreeBusyRequestItem{
			{Id: calendarID},
		},
	}).Do()

	if err != nil {
		return nil, toerr.Internal(fmt.Errorf("error GetBusyTimes, Freebusy.Query: %w", err)).Msg("error querying freebusy")
	}

	if freeBusy, ok := freeBusyResp.Calendars[calendarID]; ok {
		timeSlots := make(values.DateSlots, len(freeBusy.Busy))
		for i, b := range freeBusy.Busy {
			startTime, err := time.Parse(time.RFC3339, b.Start)
			if err != nil {
				return nil, toerr.Invalid(fmt.Errorf("error GetBusyTimes, invalid start time: %w", err)).Msg("invalid start time")
			}
			endTime, err := time.Parse(time.RFC3339, b.End)
			if err != nil {
				return nil, toerr.Invalid(fmt.Errorf("error GetBusyTimes, invalid end time: %w", err)).Msg("invalid end time")
			}
			timeSlots[i] = values.DateSlot{
				Start: startTime,
				End:   endTime,
			}
		}
		return timeSlots, nil
	}

	return values.DateSlots{}, nil
}
