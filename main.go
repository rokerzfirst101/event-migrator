package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google"
	"event-migration-script/google/calendar"
	"event-migration-script/handler/migrator"
	"event-migration-script/models"
	"os"
)

const (
	GoogleCredential    = "GOOGLE_CREDENTIALS"
	RefreshToken        = "REFRESH_TOKEN"
	SourceCalendar      = "SOURCE_CALENDAR"
	DestinationCalendar = "DESTINATION_CALENDAR"
	TimeMin             = "TIME_MIN"
	TimeMax             = "TIME_MAX"
)

func main() {
	app := gofr.NewCMD()

	app.GET("start", func(ctx *gofr.Context) (interface{}, error) {
		migratorConfig := getMigratorConfigFromEnv(ctx)

		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		eventMigrator := migrator.New(client, migratorConfig)

		pageToken, err := eventMigrator.Start(ctx)
		if err != nil {
			return nil, err
		}

		return nil, writePageTokenToFile(pageToken)
	})

	app.GET("load", func(ctx *gofr.Context) (interface{}, error) {
		migratorConfig := getMigratorConfigFromEnv(ctx)
		reverseMigratorConfig := getReverseMigratorConfigFromEnv(ctx)

		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		eventMigrator := migrator.New(client, migratorConfig)
		reverseMigrator := migrator.New(client, reverseMigratorConfig)

		for true {
			ctx.Logger.Infof("Starting event migrator")

			_, err = eventMigrator.StartWithGoroutine(ctx)
			if err != nil {
				return nil, err
			}

			ctx.Logger.Infof("Starting reverse event migrator")
			_, err = reverseMigrator.StartWithGoroutine(ctx)
			if err != nil {
				return nil, err
			}

			eventMigrator = migrator.New(client, reverseMigratorConfig)
		}

		return nil, nil
	})

	app.Start()
}

func getMigratorConfigFromEnv(ctx *gofr.Context) *models.MigratorConfig {
	return &models.MigratorConfig{
		SourceCalendarID:      ctx.Config.Get(SourceCalendar),
		DestinationCalendarID: ctx.Config.Get(DestinationCalendar),
		TimeMin:               ctx.Config.Get(TimeMin),
		TimeMax:               ctx.Config.Get(TimeMax),
	}
}

func getReverseMigratorConfigFromEnv(ctx *gofr.Context) *models.MigratorConfig {
	return &models.MigratorConfig{
		SourceCalendarID:      ctx.Config.Get(DestinationCalendar),
		DestinationCalendarID: ctx.Config.Get(SourceCalendar),
		TimeMin:               ctx.Config.Get(TimeMin),
		TimeMax:               ctx.Config.Get(TimeMax),
	}
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
