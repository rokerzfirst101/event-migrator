# Event Migration Script

The Script uses a single OAuth Account, the OAuth account being used must have write access on both the calendars.

## Usage
1. Set the env values as defined below.
2. Start the program using `go run main.go start`


## Environment Variables

```
APP_NAME: The app name to be used by Gofr
APP_VERSION: The app version 

GOOGLE_CREDENTIALS: Credentials for the GCP Project
REFRESH_TOKEN: OAuth token for the account with Calendar Scope

SOURCE_CALENDAR: Calendar ID of the source calendar
DESTINATION_CALENDAR: Calendar ID of the destination calendar
```

