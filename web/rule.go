package money

import (
	"context"
	"fmt"
	"log/slog"

	core "github.com/keeth/money/core"
	model "github.com/keeth/money/model"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetRuleExprsText(row model.GetRulesRow) []string {
	exprs := []string{}
	if row.Rule.AmountExpr.Valid && row.Rule.AmountExpr.String != "" {
		exprs = append(exprs, "amount: "+row.Rule.AmountExpr.String)
	}
	if row.Rule.DescExpr.Valid && row.Rule.DescExpr.String != "" {
		exprs = append(exprs, "desc: "+row.Rule.DescExpr.String)
	}
	if row.Rule.DateExpr.Valid && row.Rule.DateExpr.String != "" {
		exprs = append(exprs, "date: "+row.Rule.DateExpr.String)
	}
	if row.Cat.ID != 0 {
		exprs = append(exprs, "category: "+row.Cat.Name)
	}
	return exprs
}

func GetRuleRows(rows []model.GetRulesRow) Node {
	return Group{Map(rows, func(row model.GetRulesRow) Node {
		return Tr(
			If(row.Rule.ID == rows[len(rows)-1].Rule.ID,
				nextPageNode("rule", fmt.Sprintf("%d", row.Rule.Ord)),
			),
			Td(Text(row.Rule.TestExpr)),
			Td(Map(GetRuleExprsText(row), func(expr string) Node {
				return Div(Text(expr))
			})),
		)
	})}
}

func GetRules(ctx context.Context, app *core.App, params model.GetRulesParams) (error, Node) {
	rules, err := app.Model.GetRules(ctx, params)
	if err != nil {
		slog.Error("failed to get rules", "error", err)
		return err, nil
	}
	rowNodes := GetRuleRows(rules)
	if params.After == 0 {
		return tablePage(PageProps{
			Title:       "Rules",
			Description: "Rules",
		},
			[]string{"Test", "Actions"},
			rowNodes,
		)
	}
	return nil, rowNodes
}

func GetRulesEndpoint(c echo.Context) error {
	cc := c.(*Context)
	after, limit := paginationParamsInt(c)
	err, rules := GetRules(cc.Request().Context(), cc.App, model.GetRulesParams{
		After: after,
		Limit: limit,
	})
	return render(c, err, rules)
}
