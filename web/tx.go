package money

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/keeth/money"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

var maxLimit = int64(100)

func GetTxRows(txs []sqlc.GetTxsRow) Node {
	return Group{Map(txs, func(txRow sqlc.GetTxsRow) Node {
		return Tr(
			If(txRow.Tx.ID == txs[len(txs)-1].Tx.ID,
				Group{
					hx.Get(fmt.Sprintf("/tx?after=%s", url.QueryEscape(txRow.Tx.Ord))),
					hx.Trigger("revealed"),
					hx.Swap("afterend"),
				},
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
		return nil, page(PageProps{
			Title:       "Transactions",
			Description: "Transactions",
		},
			Table(
				THead(
					Tr(
						Th(Text("Date")),
						Th(Text("Description")),
						Th(Text("Amount")),
					),
				),
				TBody(
					rowNodes,
				),
			),
		)
	}
	return nil, rowNodes
}

func GetTxsEndpoint(c echo.Context) error {
	cc := c.(*Context)
	after := cc.QueryParam("after")
	limitStr := cc.QueryParam("limit")
	limit := maxLimit
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return cc.String(http.StatusBadRequest, "invalid limit")
		}
		if limit < 1 || limit > maxLimit {
			limit = maxLimit
		}
	}
	err, txs := GetTxs(cc.Request().Context(), cc.App, model.GetTxsParams{
		After: after,
		Limit: limit,
	})
	if err != nil {
		return cc.String(http.StatusInternalServerError, err.Error())
	}
	return txs.Render(cc.Response().Writer)
}
