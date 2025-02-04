package money

import (
	"context"
)

func (q *Queries) CreateOrUpdateTx_(ctx context.Context, arg CreateOrUpdateTxParams) error {
	arg.Ord = arg.Date + " " + arg.Xid
	return q.CreateOrUpdateTx(ctx, arg)
}
