package money

import (
	"context"
	"os"

	ofx "github.com/aclindsa/ofxgo"
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

func (a *App) ImportOFX(file *os.File) error {
	ofxResp, err := ofx.ParseResponse(file)
	if err != nil {
		return err
	}

	resp, err := ParseOfxResponse(ofxResp)
	if err != nil {
		return err
	}

	acc, err := a.queries.GetOrCreateAcc(a.ctx, data.CreateAccParams{
		Xid:  resp.ID,
		Kind: resp.Kind,
	}, maxAccounts)

	if err != nil {
		return err
	}

	for _, ofxTx := range resp.Transactions {
		tx := ParseOfxTransaction(&ofxTx)

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
	}
	return nil
}
