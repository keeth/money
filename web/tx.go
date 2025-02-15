package money

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/keeth/money"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetTxs(c echo.Context) error {
	app := money.GetGlobalApp()
	txs, err := app.Model.GetTxs(c.Request().Context(), model.GetTxsParams{
		Before: "9999-99-99",
		Limit:  100,
	})
	if err != nil {
		slog.Error("failed to get txs", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
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
				Map(txs, func(txRow sqlc.GetTxsRow) Node {
					return Tr(
						Td(Text(txRow.Tx.Date)),
						Td(Text(txRow.Tx.Desc)),
						Td(Text(fmt.Sprintf("%.2f", txRow.Tx.Amount))),
					)
				}),
			),
		),
	).Render(c.Response().Writer)
	return nil
}
