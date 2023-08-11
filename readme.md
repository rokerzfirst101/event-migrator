# Event Migration Script

The Script uses a single OAuth Account, the OAuth account being used must have write access on both the calendars.
By default, the code resorts to an exponential backoff of 5 seconds, this can be changed by setting the `BACKOFF_TIME` env variable
for a maximum of `MAX_RETRIES` times.

The script is setup to migrate events from first to last in chronological order.

## Usage
1. Get Dependencies using `go get ./...`
2. Create a `.local.env` file in the root directory. 
3. Set the required environment variables in the `.local.env` file.
4. Set the optional environment variables as per the requirement.
5. Start the program using `go run main.go start`

## Environment Variables

| Variable Name        | Description                                                  | Required |
|----------------------|--------------------------------------------------------------|----------|
| GOOGLE_CREDENTIALS   | Credentials for the GCP Project                              | Yes      |
| REFRESH_TOKEN        | OAuth token for the account with Calendar Scope              | Yes      |
| SOURCE_CALENDAR      | Calendar ID of the source calendar                           | Yes      |
| DESTINATION_CALENDAR | Calendar ID of the destination calendar                      | Yes      |
| EVENT_TIME_MIN       | Lower bound for the start time of the events to be migrated. | No       |
| EVENT_TIME_MAX       | Upper bound for the start time of the events to be migrated. | No       |
| MAX_RETRIES          | Maximum number of retries for the exponential backoff.       | No       |
| BACKOFF_TIME         | Time to wait between retries.                                | No       |

> Note: The `EVENT_TIME_MIN` and `EVENT_TIME_MAX` are in RFC3339 format. E.g.: `2023-08-10T21:27:37+05:30`
