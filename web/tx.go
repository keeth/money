package money

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/keeth/money"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetTxRows(txs []sqlc.GetTxsRow) Node {
	return Group{Map(txs, func(txRow sqlc.GetTxsRow) Node {
		return Tr(
			If(txRow.Tx.ID == txs[len(txs)-1].Tx.ID,
				nextPageNode("tx", txRow.Tx.Ord),
			),
			Td(Text(txRow.Tx.Date)),
			Td(Text(txRow.Tx.Desc)),
			Td(Text(fmt.Sprintf("%.2f", txRow.Tx.Amount))),
		)
	})}
}

func GetTxs(ctx context.Context, app *money.App, params model.GetTxsParams) (error, Node) {
	txs, err := app.Model.GetTxs(ctx, params)
	if err != nil {
		slog.Error("failed to get txs", "error", err)
		return err, nil
	}
	rowNodes := GetTxRows(txs)
	if params.After == "" {
		return tablePage(PageProps{
			Title:       "Transactions",
			Description: "Transactions",
		},
			[]string{"Date", "Description", "Amount"},
			rowNodes,
		)
	}
	return nil, rowNodes
}

func GetTxsEndpoint(c echo.Context) error {
	cc := c.(*Context)
	after, limit := paginationParams(c)
	err, txs := GetTxs(cc.Request().Context(), cc.App, model.GetTxsParams{
		After: after,
		Limit: limit,
	})
	return render(c, err, txs)
}
