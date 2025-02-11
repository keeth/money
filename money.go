package money

import (
	"context"
	"os"

	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
)

type App struct {
	Queries *sqlc.Queries
}

func NewApp(queries *sqlc.Queries) *App {
	return &App{
		Queries: queries,
	}
}

var maxAccounts = 10

type ImportResult struct {
	TxCreated  int
	TxUpdated  int
	AccCreated int
}

var app *App

func InitGlobalApp(q *sqlc.Queries) *App {
	if app == nil {
		app = NewApp(q)
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

	accRow, err := model.GetOrCreateAcc(ctx, a.Queries, sqlc.CreateAccParams{
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
		txRow, err := model.CreateOrUpdateTx(ctx, a.Queries, sqlc.CreateOrUpdateTxParams{
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
