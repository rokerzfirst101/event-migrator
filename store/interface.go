package store

import "developer.zopsmart.com/go/gofr/pkg/gofr"

type Interview interface {
	GetEventIDMap(ctx *gofr.Context) (map[int]string, error)
	UpdateEventID(ctx *gofr.Context, oldEventID, newEventID string) error
}
