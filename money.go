package money

import (
	"context"
	"database/sql"
	"os"

	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
)

type App struct {
	Model *model.ModelContext
}

func NewApp(db *sql.DB) *App {
	return &App{
		Model: model.NewModelContext(db),
	}
}

var maxAccounts = 10

type ImportResult struct {
	TxCreated  int
	TxUpdated  int
	AccCreated int
}

var app *App

func InitGlobalApp(db *sql.DB) *App {
	if app == nil {
		app = NewApp(db)
	}
	return app
}

func GetGlobalApp() *App {
	return app
}

func (a *App) ImportOFX(ctx context.Context, file *os.File) (ImportResult, error) {
	result := ImportResult{}

	resp, err := ParseOfxResponse(file)
	if err != nil {
		return ImportResult{}, err
	}

	accRow, err := a.Model.GetOrCreateAcc(ctx, sqlc.CreateAccParams{
		Xid:  resp.ID,
		Kind: resp.Kind,
	}, maxAccounts)
	if err != nil {
		return ImportResult{}, err
	}
	if accRow.Created {
		result.AccCreated++
	}

	for _, tx := range resp.Transactions {
		txRow, err := a.Model.CreateOrUpdateTx(ctx, sqlc.CreateOrUpdateTxParams{
			Date:   tx.Date,
			Amount: tx.Amount,
			Desc:   tx.Desc,
			AccID:  accRow.Acc.ID,
			Xid:    tx.ID,
		})
		if err != nil {
			return ImportResult{}, err
		}
		if txRow.Created {
			result.TxCreated++
		} else {
			result.TxUpdated++
		}
	}
	return result, nil
}
