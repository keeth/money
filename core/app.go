package money

import (
	"database/sql"

	model "github.com/keeth/money/model"
)

type App struct {
	Model *model.ModelContext
}

func NewApp(db *sql.DB) *App {
	return &App{
		Model: model.NewModelContext(db),
	}
}
