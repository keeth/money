package money

import (
	"context"
	"database/sql"
	"fmt"
)

func (q *Queries) CreateOrUpdateTx_(ctx context.Context, arg CreateOrUpdateTxParams) error {
	arg.Ord = arg.Date + " " + arg.Xid
	return q.CreateOrUpdateTx(ctx, arg)
}

func (q *Queries) GetOrCreateAcc(ctx context.Context, arg CreateAccParams, maxAttempts int) (Acc, error) {
	acc, err := q.GetAccByXid(ctx, arg.Xid)
	if err != nil {
		if err != sql.ErrNoRows {
			return acc, err
		}
		for i := range maxAttempts {
			var name string
			if i == 0 {
				name = arg.Kind
			} else {
				name = fmt.Sprintf("%s%d", arg.Kind, i)
			}
			err = q.CreateAcc(ctx, CreateAccParams{
				Xid:  arg.Xid,
				Kind: arg.Kind,
				Name: name,
			})
			if err == nil {
				acc, err = q.GetAccByXid(ctx, arg.Xid)
				if err != nil {
					return acc, err
				}
				break
			} else {
				if err.Error() == "UNIQUE constraint failed: acc.name" { // todo: check if this is the correct error
					continue
				}
				return acc, err
			}
		}
		if acc.ID == 0 {
			return acc, fmt.Errorf("failed to create account after %d attempts", maxAttempts)
		}
	}
	return acc, err
}
