package helpers

import (
	"context"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetContextWithMockDB() (*gofr.Context, sqlmock.Sqlmock) {
	app := gofr.New()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	app.ORM = nil

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})

	app.SetORM(datastore.GORMClient{DB: gormDB})

	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	return ctx, mock
}
