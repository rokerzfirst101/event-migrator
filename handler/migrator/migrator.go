package migrator

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/handler"
	"event-migration-script/models"
	gcalendar "google.golang.org/api/calendar/v3"
)

type migrator struct {
	calendar google.CalendarService
	config   *models.MigratorConfig
}

func New(client google.CalendarService, config *models.MigratorConfig) handler.EventMigrator {
	return &migrator{calendar: client, config: config}
}

func (m *migrator) Start(ctx *gofr.Context) (interface{}, error) {
	pageSyncToken := ""
	firstRun := true

	for firstRun || pageSyncToken != "" {
		events, newPageSyncToken, err := m.calendar.List(m.config.SourceCalendarID, pageSyncToken, m.config.TimeMin, m.config.TimeMax)
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
		if events[i].Organizer.Email == m.config.SourceCalendarID {
			_, err := m.calendar.Move(m.config.SourceCalendarID, events[i].Id, m.config.DestinationCalendarID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
