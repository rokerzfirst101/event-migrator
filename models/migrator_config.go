package models

type MigratorConfig struct {
	SourceCalendarID      string `json:"sourceCalendarID"`
	DestinationCalendarID string `json:"destinationCalendarID"`
	TimeMin               string `json:"timeMin"`
	TimeMax               string `json:"timeMax"`
}
