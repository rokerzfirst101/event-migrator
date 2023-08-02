package migrator

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/handler"
	"event-migration-script/store"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
)

type migrator struct {
	interview   store.Interview
	oldCalendar google.CalendarService
	newCalendar google.CalendarService
}

func New(interview store.Interview, oldCalendar google.CalendarService, newCalendar google.CalendarService) handler.EventMigrator {
	return &migrator{interview: interview, oldCalendar: oldCalendar, newCalendar: newCalendar}
}

func (m *migrator) Start(ctx *gofr.Context) (interface{}, error) {
	eventIDMap, err := m.interview.GetEventIDMap(ctx)
	if err != nil {
		return nil, err
	}

	for interviewID, eventID := range eventIDMap {
		err = m.migrateEvent(ctx, interviewID, eventID)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (m *migrator) migrateEvent(ctx *gofr.Context, interviewID int, eventID string) error {
	event, err := m.oldCalendar.Get(eventID)
	if v, ok := err.(*googleapi.Error); ok && v.Code == 410 {
		return m.interview.UpdateEventID(ctx, eventID, "")
	}

	sanitiseEvent(event)

	newEvent, err := m.newCalendar.Create(event)
	if err != nil {
		return err
	}

	err = m.interview.UpdateEventID(ctx, eventID, newEvent.Id)
	if err != nil {
		return err
	}

	return nil
}

func sanitiseEvent(event *calendar.Event) {
	event.Id, event.HtmlLink, event.ICalUID, event.Etag = "", "", "", ""
}
