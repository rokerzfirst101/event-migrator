package migrator

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/handler"
	"event-migration-script/models"
	"event-migration-script/pretty"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"os"
	"sync"
	"time"
)

const (
	MaxRetries        = 5
	IntervalIncrement = 5 * time.Second
)

type migrator struct {
	calendar google.CalendarService
	config   *models.MigratorConfig
}

func New(client google.CalendarService, config *models.MigratorConfig) handler.EventMigrator {
	return &migrator{calendar: client, config: config}
}

func (m *migrator) Start(ctx *gofr.Context) (interface{}, error) {
	pageSyncToken := "START"
	count := 0

	for pageSyncToken != "" && count < MaxRetries {
		ctx.Logger.Infof("Calling List API for pageSyncToken: %s", pageSyncToken)

		if pageSyncToken == "START" {
			pageSyncToken = ""
		}

		events, newPageSyncToken, err := m.calendar.List(m.config, 100, pageSyncToken)
		if err != nil {
			googleErr, _ := err.(*googleapi.Error)
			if googleErr.Code == 404 {
				ctx.Logger.Infof("Calendar not found. Skipping migration")

				return nil, err
			} else if googleErr.Code == 403 {
				sleepTime := time.Duration(count) * IntervalIncrement

				ctx.Logger.Infof("Quota Exceeded. Sleeping for %v", sleepTime)

				time.Sleep(sleepTime)

				count++

				continue
			} else {
				return nil, err
			}
		}

		for _, event := range events {
			err = m.migrateEvent(ctx, event)
			if err != nil {
				logFailedIDToFile(ctx, event.Id)

				continue
			}

			count++
		}

		pageSyncToken = newPageSyncToken
	}

	return pageSyncToken, nil
}

func logFailedIDToFile(ctx *gofr.Context, eventID string) {
	ctx.Logger.Infof("Logging failed event ID to file failed_ids.txt")

	file, _ := os.OpenFile("failed_ids.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	_, _ = file.WriteString(eventID + "\n")
}

func (m *migrator) StartWithGoroutine(ctx *gofr.Context) (interface{}, error) {
	wg := sync.WaitGroup{}

	pageSyncToken := ""
	firstRun := true

	for firstRun || pageSyncToken != "" {
		ctx.Logger.Infof("Calling List API for pageSyncToken: %s", pageSyncToken)

		events, newPageSyncToken, err := m.calendar.List(m.config, 100, pageSyncToken)
		if err != nil {
			ctx.Logger.Errorf("Sleeping: %v", err)

			time.Sleep(5 * time.Second)

			continue
		}

		for i := range events {
			i := i

			wg.Add(1)
			go func() {
				defer wg.Done()

				err := m.migrateEvent(ctx, events[i])
				if err != nil {
					ctx.Logger.Errorf("Migrate Event Failed, Returning: %v", err)

					return
				}
			}()
		}

		firstRun = false
		pageSyncToken = newPageSyncToken
	}

	wg.Wait()

	return pageSyncToken, nil
}

func (m *migrator) migrateEvent(ctx *gofr.Context, event *gcalendar.Event) error {
	ctx.Logger.Infof("Starting Migration for Event: %s, from: %s to %s", pretty.PrintEvent(event), m.config.SourceCalendarID,
		m.config.DestinationCalendarID)
	if event.Organizer.Email == m.config.SourceCalendarID {
		maxRetries := 0

		for maxRetries < 3 {
			_, err := m.calendar.Move(m.config.SourceCalendarID, event.Id, m.config.DestinationCalendarID)
			if err != nil {
				sleepTime := time.Duration(maxRetries+1) * IntervalIncrement

				ctx.Logger.Errorf("Got an error sleeping for %v: %v", sleepTime, err)

				time.Sleep(sleepTime)

				maxRetries++

				continue
			}

			break
		}
	}

	ctx.Logger.Infof("Completed Migration for Event: %s, from: %s to %s", pretty.PrintEvent(event), m.config.SourceCalendarID,
		m.config.DestinationCalendarID)

	return nil
}
