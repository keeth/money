package money

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateOrUpdateTxResult struct {
	ID      int64
	Created bool
}

func (q *Queries) CreateOrUpdateTx_(ctx context.Context, arg CreateOrUpdateTxParams) (CreateOrUpdateTxResult, error) {
	arg.Ord = arg.Date + " " + arg.Xid
	row, err := q.CreateOrUpdateTx(ctx, arg)
	result := CreateOrUpdateTxResult{}
	if err != nil {
		return result, err
	}
	result.ID = row.ID
	result.Created = row.CreatedAt == row.UpdatedAt
	return result, nil
}

type GetOrCreateAccResult struct {
	Acc     Acc
	Created bool
}

func (q *Queries) GetOrCreateAcc(ctx context.Context, arg CreateAccParams, maxAttempts int) (GetOrCreateAccResult, error) {
	result := GetOrCreateAccResult{}
	acc, err := q.GetAccByXid(ctx, arg.Xid)
	if err != nil {
		if err != sql.ErrNoRows {
			return result, err
		}
		for i := range maxAttempts {
			var name string
			if i == 0 {
				name = arg.Kind
			} else {
				name = fmt.Sprintf("%s%d", arg.Kind, i)
			}
			_, err = q.CreateAcc(ctx, CreateAccParams{
				Xid:  arg.Xid,
				Kind: arg.Kind,
				Name: name,
			})
			if err == nil {
				result.Created = true
				acc, err = q.GetAccByXid(ctx, arg.Xid)
				if err != nil {
					return result, err
				}
				break
			} else {
				if err.Error() == "UNIQUE constraint failed: acc.name" { // todo: check if this is the correct error
					continue
				}
				return result, err
			}

		}
		if acc.ID == 0 {
			return result, fmt.Errorf("failed to create account after %d attempts", maxAttempts)
		}
	}
	result.Acc = acc
	return result, nil
}
