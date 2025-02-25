package money

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
)

var exprCache = make(map[string]*vm.Program)

func compileExpr(code string) (*vm.Program, error) {
	if program, ok := exprCache[code]; ok {
		return program, nil
	}
	program, err := expr.Compile(code, expr.Env(sqlc.Tx{}))
	if err != nil {
		return nil, err
	}
	exprCache[code] = program
	return program, nil
}

func compileTest(rule model.Rule) (*vm.Program, error) {
	return compileExpr(rule.TestExpr)
}

func evaluateTest(test *vm.Program, tx sqlc.Tx) (bool, error) {
	eval, err := expr.Run(test, tx)
	if err != nil {
		return false, err
	}
	testResult, ok := eval.(bool)
	if !ok {
		return false, fmt.Errorf("test expression must return a boolean")
	}
	return testResult, nil
}

func compileTests(rules []model.Rule) ([]*vm.Program, error) {
	tests := []*vm.Program{}
	for _, rule := range rules {
		test, err := compileTest(rule)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return tests, nil
}

func processExpr[T any](exprStr string, tx sqlc.Tx) (T, error) {
	var zero T
	program, err := compileExpr(exprStr)
	if err != nil {
		return zero, err
	}
	result, err := expr.Run(program, tx)
	if err != nil {
		return zero, err
	}
	newValue, ok := result.(T)
	if !ok {
		return zero, fmt.Errorf("expression must return a value of type %T", zero)
	}
	return newValue, nil
}

func ApplyRules(ctx context.Context, rules []model.Rule, tx *sqlc.Tx) (bool, error) {
	tests, err := compileTests(rules)
	if err != nil {
		return false, err
	}
	// save the currently persisted values of the tx for later comparison
	prevAmount := tx.Amount
	prevDate := tx.Date
	prevDesc := tx.Desc
	prevCatID := tx.CatID

	// reset the tx to any original values before applying the rules.
	// rule application should be idempotent.
	if tx.OrigAmount.Valid {
		tx.Amount = tx.OrigAmount.Float64
	}
	if tx.OrigDate.Valid {
		tx.Date = tx.OrigDate.String
	}
	if tx.OrigDesc.Valid {
		tx.Desc = tx.OrigDesc.String
	}
	tx.CatID = sql.NullInt64{}

	// if rule matches, apply the changes
	for i, test := range tests {
		testResult, err := evaluateTest(test, *tx)
		if err != nil {
			return false, err
		}
		if !testResult {
			continue
		}
		rule := rules[i]

		if rule.CatID.Valid {
			tx.CatID = rule.CatID
		}

		if rule.AmountExpr.Valid {
			if amount, err := processExpr[float64](rule.AmountExpr.String, *tx); err != nil {
				return false, err
			} else {
				tx.Amount = amount
			}
		}

		if rule.DescExpr.Valid {
			if desc, err := processExpr[string](rule.DescExpr.String, *tx); err != nil {
				return false, err
			} else {
				tx.Desc = desc
			}
		}

		if rule.DateExpr.Valid {
			if date, err := processExpr[string](rule.DateExpr.String, *tx); err != nil {
				return false, err
			} else {
				tx.Date = date
			}
		}
	}
	// check whether the rules made any changes to the tx, from what was previously persisted.
	// if a rule changes a field, start preserving the original value in a separate field (copy on write).
	dirty := false
	if prevCatID != tx.CatID {
		dirty = true
	}
	if prevAmount != tx.Amount {
		if !tx.OrigAmount.Valid {
			tx.OrigAmount = sql.NullFloat64{
				Float64: prevAmount,
				Valid:   true,
			}
		}
		dirty = true
	}
	if prevDesc != tx.Desc {
		if !tx.OrigDesc.Valid {
			tx.OrigDesc = sql.NullString{
				String: prevDesc,
				Valid:  true,
			}
		}
		dirty = true
	}
	if prevDate != tx.Date {
		if !tx.OrigDate.Valid {
			tx.OrigDate = sql.NullString{
				String: prevDate,
				Valid:  true,
			}
		}
		dirty = true
	}
	return dirty, nil
}
