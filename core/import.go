package money

import (
	"context"
	"os"

	sqlc "github.com/keeth/money/model/sqlc"
)

var maxAccounts = 10

type ImportResult struct {
	TxCreated  int
	TxUpdated  int
	AccCreated int
}

func (a *App) ImportOFX(ctx context.Context, file *os.File) (ImportResult, error) {
	result := ImportResult{}

	resp, err := ParseOfxResponse(file)
	if err != nil {
		return result, err
	}

	accRow, err := a.Model.GetOrCreateAcc(ctx, sqlc.CreateAccParams{
		Xid:  resp.ID,
		Kind: resp.Kind,
	}, maxAccounts)
	if err != nil {
		return result, err
	}
	if accRow.Created {
		result.AccCreated++
	}
	rules, err := a.Model.GetAllRules(ctx)
	if err != nil {
		return result, err
	}
	for _, tx := range resp.Transactions {
		_, err := ApplyRules(ctx, rules, &tx)
		if err != nil {
			return result, err
		}
		txRow, err := a.Model.CreateOrUpdateTx(ctx, sqlc.CreateOrUpdateTxParams{
			Date:       tx.Date,
			OrigDate:   tx.OrigDate,
			Amount:     tx.Amount,
			OrigAmount: tx.OrigAmount,
			Desc:       tx.Desc,
			OrigDesc:   tx.OrigDesc,
			AccID:      accRow.Acc.ID,
			Xid:        tx.Xid,
			CatID:      tx.CatID,
		})
		if err != nil {
			return result, err
		}
		if txRow.Created {
			result.TxCreated++
		} else {
			result.TxUpdated++
		}
	}
	return result, nil
}
