package money

import (
	"context"
	"os"

	data "github.com/keeth/money/data"
)

type App struct {
	ctx     context.Context
	queries *data.Queries
}

func NewApp(ctx context.Context, q *data.Queries) *App {
	return &App{
		ctx:     ctx,
		queries: q,
	}
}

var maxAccounts = 10

type ImportResult struct {
	TxCreated  int
	TxUpdated  int
	AccCreated int
}

func (a *App) ImportOFX(file *os.File) (ImportResult, error) {
	result := ImportResult{}

	resp, err := ParseOfxResponse(file)
	if err != nil {
		return ImportResult{}, err
	}

	acc, err := a.queries.GetOrCreateAcc(a.ctx, data.CreateAccParams{
		Xid:  resp.ID,
		Kind: resp.Kind,
	}, maxAccounts)
	if err != nil {
		return ImportResult{}, err
	}

	for _, tx := range resp.Transactions {
		err = a.queries.CreateOrUpdateTx_(a.ctx, data.CreateOrUpdateTxParams{
			Date:       tx.Date,
			Amount:     tx.Amount,
			Desc:       tx.Desc,
			AccID:      acc.ID,
			Xid:        tx.ID,
			OrigDate:   tx.Date,
			OrigAmount: tx.Amount,
			OrigDesc:   tx.Desc,
		})
		if err != nil {
			return ImportResult{}, err
		}
		result.TxCreated++
	}
	return result, nil
}
