package money

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
)

func compileExpr(code string) (*vm.Program, error) {
	return expr.Compile(code, expr.Env(sqlc.Tx{}))
}

func compileTest(rule model.Rule) (*vm.Program, error) {
	return compileExpr(rule.TestExpr)
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

func getRulesFromRows(rows []model.GetRulesRow) []model.Rule {
	rules := []model.Rule{}
	for _, row := range rows {
		rules = append(rules, row.Rule)
	}
	return rules
}

func ApplyRules(ctx context.Context, mc model.ModelContext, tx sqlc.Tx) error {
	rows, err := mc.GetRules(ctx, model.GetRulesParams{
		After: 0,
		Limit: 100,
	})
	if err != nil {
		return err
	}
	rules := getRulesFromRows(rows)
	tests, err := compileTests(rules)
	if err != nil {
		return err
	}
	for _, test := range tests {
		eval, err := expr.Run(test, tx)
		if err != nil {
			return err
		}
		fmt.Println(eval)
	}
	return nil
}
