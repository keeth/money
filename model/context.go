package money

import (
	"database/sql"

	sqlc "github.com/keeth/money/model/sqlc"
)

type ModelContext struct {
	Queries *sqlc.Queries
	DB      *sql.DB
}

func NewModelContext(sqldb *sql.DB) *ModelContext {
	return &ModelContext{
		Queries: sqlc.New(sqldb),
		DB:      sqldb,
	}
}
