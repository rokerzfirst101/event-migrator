package pretty

import (
	"fmt"
	"google.golang.org/api/calendar/v3"
)

func PrintEvent(event *calendar.Event) string {
	return fmt.Sprintf("ID: %s, Title: %s, StartTime: %s", event.Id, event.Summary, event.Start.DateTime)
}
