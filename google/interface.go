package google

import (
	"event-migration-script/models"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"google.golang.org/api/calendar/v3"
)

//go:generate mockgen -source=interface.go -destination=mock_interface.go -package=google

type ClientProvider interface {
	GetClient(ctx *gofr.Context, state string, scopes []string) (*http.Client, error)
}

type CalendarService interface {
	List(m *models.MigratorConfig, batchSize int, pageToken string) ([]*calendar.Event, string, error)
	Move(sourceCalendar, eventID, destinationCalendarID string) (*calendar.Event, error)
}
