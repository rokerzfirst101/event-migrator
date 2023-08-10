package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/google/calendar"
	"event-migration-script/handler/migrator"
)

const (
	GoogleCredential    = "GOOGLE_CREDENTIALS"
	RefreshToken        = "REFRESH_TOKEN"
	SourceCalendar      = "SOURCE_CALENDAR"
	DestinationCalendar = "DESTINATION_CALENDAR"
)

func main() {
	app := gofr.NewCMD()

	app.GET("start", func(ctx *gofr.Context) (interface{}, error) {
		sourceCalendarID := ctx.Config.Get(SourceCalendar)
		destinationCalendarID := ctx.Config.Get(DestinationCalendar)

		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		eventMigrator := migrator.New(client, sourceCalendarID, destinationCalendarID)

		_, err = eventMigrator.Start(ctx)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	app.Start()
}

func getClient(ctx *gofr.Context) (google.CalendarService, error) {
	client, err := calendar.New(ctx, GoogleCredential, RefreshToken)
	if err != nil {
		return nil, err
	}

	return client, nil
}
