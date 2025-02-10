package money

import (
	"context"
	"os"

	data "github.com/keeth/money/data"
)

type App struct {
	Queries *data.Queries
}

func NewApp(queries *data.Queries) *App {
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

func InitGlobalApp(q *data.Queries) *App {
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

	accRow, err := a.Queries.GetOrCreateAcc(ctx, data.CreateAccParams{
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
		txRow, err := a.Queries.CreateOrUpdateTx_(ctx, data.CreateOrUpdateTxParams{
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
