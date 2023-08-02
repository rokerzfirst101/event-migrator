package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/google/calendar"
	"event-migration-script/handler"
	"event-migration-script/handler/migrator"
	"event-migration-script/store/interview"
)

const (
	OldProjectCredentialsEnv = "OLD_PROJECT_CREDENTIALS"
	OldRefreshTokensEnv      = "OLD_REFRESH_TOKENS"
	NewProjectCredentialsEnv = "NEW_PROJECT_CREDENTIALS"
	NewRefreshTokensEnv      = "NEW_REFRESH_TOKENS"
)

func main() {
	app := gofr.NewCMD()

	app.GET("start", func(c *gofr.Context) (interface{}, error) {
		eventMigrationHandler, err := getEventMigrationHandler(c)
		if err != nil {
			return nil, err
		}

		_, err = eventMigrationHandler.Start(c)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	app.Start()
}

func getEventMigrationHandler(ctx *gofr.Context) (handler.EventMigrator, error) {
	interviewStore := interview.New()

	oldClient, err := calendar.New(ctx, OldProjectCredentialsEnv, OldRefreshTokensEnv, "")
	if err != nil {
		return nil, err
	}

	newClient, err := calendar.New(ctx, OldProjectCredentialsEnv, OldRefreshTokensEnv, "")
	if err != nil {
		return nil, err
	}

	return migrator.New(interviewStore, oldClient, newClient), nil
}
