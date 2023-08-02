package calendar

import (
	"context"
	"event-migration-script/google"
	"event-migration-script/google/token"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const state = "calendar"

type calendar struct {
	calendarID string
	service    *gcalendar.EventsService
}

func New(ctx *gofr.Context, credentialEnvKey, tokenEnvKey, calendarID string) (google.CalendarService, error) {
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

	return &calendar{service: service.Events, calendarID: calendarID}, nil
}

func (c calendar) Get(eventID string) (*gcalendar.Event, error) {
	event, err := c.service.Get(c.calendarID, eventID).Do()
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c calendar) Create(event *gcalendar.Event) (*gcalendar.Event, error) {
	createCall := c.service.Insert(c.calendarID, event).
		SupportsAttachments(true).
		ConferenceDataVersion(1).
		SendUpdates("none")

	createdEvent, err := createCall.Do()
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

func (c calendar) Delete(eventID, sendUpdates string) error {
	return c.service.Delete(c.calendarID, eventID).SendUpdates(sendUpdates).Do()
}
