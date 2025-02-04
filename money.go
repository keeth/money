package money

import (
	"context"
	"fmt"
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

func (a *App) ImportQFX(acc *data.Acc, file *os.File) error {
	resp, err := ofx.ParseResponse(file)
	if err != nil {
		return err
	}

	var ofxTransactions []ofx.Transaction

	if len(resp.Bank) > 0 {
		if stmt, ok := resp.Bank[0].(*ofx.StatementResponse); ok {
			ofxTransactions = stmt.BankTranList.Transactions
		}
	} else if len(resp.CreditCard) > 0 {
		if stmt, ok := resp.CreditCard[0].(*ofx.StatementResponse); ok {
			ofxTransactions = stmt.BankTranList.Transactions
		}
	} else {
		return fmt.Errorf("no information found in file")
	}

	for _, ofxTx := range ofxTransactions {
		tx := ParseOfxTransaction(&ofxTx)

		err = a.queries.CreateOrUpdateTx_(a.ctx, data.CreateOrUpdateTxParams{
			Date:       tx.Date,
			Amount:     tx.Amount,
			Desc:       tx.Desc,
			AccID:      acc.ID,
			Ord:        "",
			Xid:        tx.ID,
			OrigDate:   tx.Date,
			OrigAmount: tx.Amount,
			OrigDesc:   tx.Desc,
		})
	}
	return nil
}
