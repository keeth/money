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

	accRow, err := a.queries.GetOrCreateAcc(a.ctx, data.CreateAccParams{
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
		txRow, err := a.queries.CreateOrUpdateTx_(a.ctx, data.CreateOrUpdateTxParams{
			Date:       tx.Date,
			Amount:     tx.Amount,
			Desc:       tx.Desc,
			AccID:      accRow.Acc.ID,
			Xid:        tx.ID,
			OrigDate:   tx.Date,
			OrigAmount: tx.Amount,
			OrigDesc:   tx.Desc,
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
