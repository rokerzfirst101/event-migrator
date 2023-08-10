package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/google/calendar"
	"event-migration-script/handler/migrator"
	"os"
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

		pageToken, err := eventMigrator.Start(ctx)
		if err != nil {
			return nil, err
		}

		return nil, writePageTokenToFile(pageToken)
	})

	app.Start()
}

func writePageTokenToFile(pageToken interface{}) error {
	pageTokenString := pageToken.(string)

	if pageTokenString != "" {
		file, _ := os.OpenFile("page_token.txt", os.O_CREATE|os.O_WRONLY, 0644)

		_, err := file.WriteString(pageTokenString)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func getClient(ctx *gofr.Context) (google.CalendarService, error) {
	client, err := calendar.New(ctx, GoogleCredential, RefreshToken)
	if err != nil {
		return nil, err
	}

	return client, nil
}
