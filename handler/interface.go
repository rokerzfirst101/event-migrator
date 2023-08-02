package handler

import "developer.zopsmart.com/go/gofr/pkg/gofr"

type EventMigrator interface {
	Start(ctx *gofr.Context) (interface{}, error)
}
