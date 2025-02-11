package money

import (
	"database/sql"

	sqlc "github.com/keeth/money/model/sqlc"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

type ModelContext struct {
	Queries *sqlc.Queries
	DB      *bun.DB
}

func NewModelContext(sqldb *sql.DB) *ModelContext {
	return &ModelContext{
		Queries: sqlc.New(sqldb),
		DB:      bun.NewDB(sqldb, sqlitedialect.New()),
	}
}
