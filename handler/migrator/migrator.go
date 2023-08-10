package migrator

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/handler"
	gcalendar "google.golang.org/api/calendar/v3"
)

type migrator struct {
	calendar              google.CalendarService
	sourceCalendarID      string
	destinationCalendarID string
}

func New(client google.CalendarService, sourceCalendarID, destinationCalendarID string) handler.EventMigrator {
	return &migrator{calendar: client, sourceCalendarID: sourceCalendarID, destinationCalendarID: destinationCalendarID}
}

func (m *migrator) Start(ctx *gofr.Context) (interface{}, error) {
	pageSyncToken := ""
	firstRun := true

	for firstRun || pageSyncToken != "" {
		events, newPageSyncToken, err := m.calendar.List(m.sourceCalendarID, pageSyncToken)
		if err != nil {
			continue
		}

		err = m.migrateEvents(events)
		if err != nil {
			continue
		}

		firstRun = false
		pageSyncToken = newPageSyncToken
	}

	return pageSyncToken, nil
}

func (m *migrator) migrateEvents(events []*gcalendar.Event) error {
	for i := range events {
		if events[i].Organizer.Email == m.sourceCalendarID {
			_, err := m.calendar.Move(m.sourceCalendarID, events[i].Id, m.destinationCalendarID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
