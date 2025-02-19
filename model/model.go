package money

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	sqlc "github.com/keeth/money/model/sqlc"
)

type CreateOrUpdateTxResult struct {
	ID      int64
	Created bool
}

func (mc *ModelContext) CreateOrUpdateTx(ctx context.Context, arg sqlc.CreateOrUpdateTxParams) (CreateOrUpdateTxResult, error) {
	arg.Ord = arg.Date + " " + arg.Xid
	row, err := mc.Queries.CreateOrUpdateTx(ctx, arg)
	result := CreateOrUpdateTxResult{}
	if err != nil {
		return result, err
	}
	result.ID = row.ID
	result.Created = row.CreatedAt == row.UpdatedAt
	return result, nil
}

type GetOrCreateAccResult struct {
	Acc     sqlc.Acc
	Created bool
}

func (mc *ModelContext) GetOrCreateAcc(ctx context.Context, arg sqlc.CreateAccParams, maxAttempts int) (GetOrCreateAccResult, error) {
	result := GetOrCreateAccResult{}
	acc, err := mc.Queries.GetAccByXid(ctx, arg.Xid)
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
			_, err = mc.Queries.CreateAcc(ctx, sqlc.CreateAccParams{
				Xid:  arg.Xid,
				Kind: arg.Kind,
				Name: name,
			})
			if err == nil {
				result.Created = true
				acc, err = mc.Queries.GetAccByXid(ctx, arg.Xid)
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

type GetTxsParams struct {
	After string
	Limit int64
}

func (mc *ModelContext) GetTxs(ctx context.Context, arg GetTxsParams) ([]sqlc.GetTxsRow, error) {
	if arg.Limit == 0 {
		arg.Limit = 100
	}
	if arg.After == "" {
		arg.After = "9999-12-31"
	}
	stmt := sq.Select().
		Column("tx.id").
		Column("tx.created_at").
		Column("tx.updated_at").
		Column("tx.xid").
		Column("tx.date").
		Column("tx.desc").
		Column("tx.amount").
		Column("tx.orig_date").
		Column("tx.orig_desc").
		Column("tx.orig_amount").
		Column("tx.ord").
		Column("acc.id").
		Column("acc.created_at").
		Column("acc.updated_at").
		Column("acc.name").
		Column("acc.xid").
		Column("acc.kind").
		Column("acc.is_active").
		Column("COALESCE(cat.id, 0)").
		Column("COALESCE(cat.created_at, '')").
		Column("COALESCE(cat.updated_at, '')").
		Column("COALESCE(cat.name, '')").
		Column("COALESCE(cat.kind, '')").
		Column("COALESCE(cat.is_active, 0)").
		From("tx").
		Join("acc ON tx.acc_id = acc.id").
		LeftJoin("cat ON tx.cat_id = cat.id").
		Where(sq.Lt{"tx.ord": arg.After}).
		OrderBy("tx.ord DESC").
		Limit(uint64(arg.Limit))
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	// remainder of this function is copied from sqlc/query.sql.go
	rows, err := mc.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sqlc.GetTxsRow
	for rows.Next() {
		var i sqlc.GetTxsRow
		if err := rows.Scan(
			&i.Tx.ID,
			&i.Tx.CreatedAt,
			&i.Tx.UpdatedAt,
			&i.Tx.Xid,
			&i.Tx.Date,
			&i.Tx.Desc,
			&i.Tx.Amount,
			&i.Tx.OrigDate,
			&i.Tx.OrigDesc,
			&i.Tx.OrigAmount,
			&i.Tx.Ord,
			&i.Acc.ID,
			&i.Acc.CreatedAt,
			&i.Acc.UpdatedAt,
			&i.Acc.Name,
			&i.Acc.Xid,
			&i.Acc.Kind,
			&i.Acc.IsActive,
			&i.Cat.ID,
			&i.Cat.CreatedAt,
			&i.Cat.UpdatedAt,
			&i.Cat.Name,
			&i.Cat.Kind,
			&i.Cat.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type GetCatsParams struct {
	After string
	Limit int64
}

type Cat = sqlc.Cat

func (mc *ModelContext) GetCats(ctx context.Context, arg GetCatsParams) ([]Cat, error) {
	if arg.Limit == 0 {
		arg.Limit = 100
	}
	stmt := sq.Select().
		Column("cat.id").
		Column("cat.created_at").
		Column("cat.updated_at").
		Column("cat.name").
		Column("cat.kind").
		Column("cat.is_active").
		From("cat").
		OrderBy("cat.name").
		Where(sq.Gt{"cat.name": arg.After}).
		Limit(uint64(arg.Limit))
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	// remainder of this function is copied from sqlc/query.sql.go
	rows, err := mc.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Cat
	for rows.Next() {
		var i Cat
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Kind,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type Rule = sqlc.Rule
type GetRulesRow sqlc.GetRulesRow

type GetRulesParams struct {
	After int64
	Limit int64
}

func (mc *ModelContext) GetRules(ctx context.Context, arg GetRulesParams) ([]GetRulesRow, error) {
	if arg.Limit == 0 {
		arg.Limit = 100
	}
	stmt := sq.Select().
		Column("rule.id").
		Column("rule.created_at").
		Column("rule.updated_at").
		Column("rule.start_date").
		Column("rule.end_date").
		Column("rule.test_expr").
		Column("rule.cat_id").
		Column("rule.amount_expr").
		Column("rule.desc_expr").
		Column("rule.date_expr").
		Column("rule.ord").
		Column("COALESCE(cat.id, 0)").
		Column("COALESCE(cat.created_at, '')").
		Column("COALESCE(cat.updated_at, '')").
		Column("COALESCE(cat.name, '')").
		Column("COALESCE(cat.kind, '')").
		Column("COALESCE(cat.is_active, 0)").
		From("rule").
		LeftJoin("cat ON rule.cat_id = cat.id").
		OrderBy("rule.ord").
		Where(sq.Gt{"rule.ord": arg.After}).
		Limit(uint64(arg.Limit))
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	// remainder of this function is copied from sqlc/query.sql.go
	rows, err := mc.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRulesRow
	for rows.Next() {
		var i GetRulesRow
		if err := rows.Scan(
			&i.Rule.ID,
			&i.Rule.CreatedAt,
			&i.Rule.UpdatedAt,
			&i.Rule.StartDate,
			&i.Rule.EndDate,
			&i.Rule.TestExpr,
			&i.Rule.CatID,
			&i.Rule.AmountExpr,
			&i.Rule.DescExpr,
			&i.Rule.DateExpr,
			&i.Rule.Ord,
			&i.Cat.ID,
			&i.Cat.CreatedAt,
			&i.Cat.UpdatedAt,
			&i.Cat.Name,
			&i.Cat.Kind,
			&i.Cat.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (mc *ModelContext) CreateRule(ctx context.Context, arg sqlc.CreateRuleParams) (int64, error) {
	id, err := mc.Queries.CreateRule(ctx, arg)
	if err != nil {
		return 0, err
	}
	// default ordering key is the rule ID
	err = mc.Queries.UpdateRuleOrd(ctx, sqlc.UpdateRuleOrdParams{
		ID:  id,
		Ord: id,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
