package models

import "time"

type MigratorConfig struct {
	SourceCalendarID      string        `json:"sourceCalendarID"`
	DestinationCalendarID string        `json:"destinationCalendarID"`
	TimeMin               string        `json:"timeMin"`
	TimeMax               string        `json:"timeMax"`
	MaxRetries            int           `json:"maxRetries"`
	BackoffTime           time.Duration `json:"backoffTime"`
}
