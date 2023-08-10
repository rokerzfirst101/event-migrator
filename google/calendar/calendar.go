package calendar

import (
	"context"
	"event-migration-script/google"
	"event-migration-script/google/token"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type calendar struct {
	service *gcalendar.EventsService
}

func New(ctx *gofr.Context, credentialEnvKey, tokenEnvKey string) (google.CalendarService, error) {
	client, err := token.NewClientFromEnv(ctx, credentialEnvKey, tokenEnvKey, []string{gcalendar.CalendarScope,
		gcalendar.CalendarReadonlyScope, gcalendar.CalendarEventsScope, gcalendar.CalendarEventsReadonlyScope})
	if err != nil {
		return nil, err
	}

	clientOption := option.WithHTTPClient(client)

	service, err := gcalendar.NewService(context.Background(), clientOption)
	if err != nil {
		return nil, err
	}

	return &calendar{service: service.Events}, nil
}

func (c calendar) List(sourceCalendar string, pageToken, timeMin, timeMax string) ([]*gcalendar.Event, string, error) {
	listCall := c.service.List(sourceCalendar).
		SingleEvents(true).
		OrderBy("startTime").
		MaxResults(100).
		PageToken(pageToken)

	if timeMin != "" {
		listCall.TimeMin(timeMin)
	}

	if timeMax != "" {
		listCall.TimeMax(timeMax)
	}

	events, err := listCall.Do()
	if err != nil {
		return nil, "", err
	}

	return events.Items, events.NextPageToken, nil
}

func (c calendar) Move(sourceCalendar, eventID, destinationCalendarID string) (*gcalendar.Event, error) {
	moveCall := c.service.Move(sourceCalendar, eventID, destinationCalendarID).
		SendUpdates("none")

	movedEvent, err := moveCall.Do()
	if err != nil {
		return nil, err
	}

	return movedEvent, nil
}
