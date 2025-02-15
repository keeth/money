package money

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/keeth/money"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

var limit = int64(100)

func txRowNodes(txs []sqlc.GetTxsRow) Node {
	return Group{Map(txs, func(txRow sqlc.GetTxsRow) Node {
		return Tr(
			If(txRow.Tx.ID == txs[len(txs)-1].Tx.ID,
				Group{
					hx.Get(fmt.Sprintf("/tx?before=%s", url.QueryEscape(txRow.Tx.Ord))),
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

func GetTxs(c echo.Context) error {
	app := money.GetGlobalApp()
	before := c.QueryParam("before")
	txs, err := app.Model.GetTxs(c.Request().Context(), model.GetTxsParams{
		Before: before,
		Limit:  limit,
	})
	if err != nil {
		slog.Error("failed to get txs", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	rowNodes := txRowNodes(txs)
	if before == "" {
		page(PageProps{
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
		).Render(c.Response().Writer)
	} else {
		rowNodes.Render(c.Response().Writer)
	}

	return nil
}

func GetTxsNext(c echo.Context) error {
	app := money.GetGlobalApp()
	before := c.QueryParam("before")
	txs, err := app.Model.GetTxs(c.Request().Context(), model.GetTxsParams{
		Before: before,
		Limit:  limit,
	})
	if err != nil {
		slog.Error("failed to get tx", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	txRowNodes(txs).Render(c.Response().Writer)
	return nil
}
