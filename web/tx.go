package money

import (
	"fmt"
	"net/http"

	"github.com/keeth/money"
	data "github.com/keeth/money/data"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetTxs(c echo.Context) error {
	app := money.GetGlobalApp()
	txs, err := app.Queries.GetTxs(c.Request().Context(), data.GetTxsParams{
		Ord:   "9999-99-99",
		Limit: 10,
	})
	if err != nil {
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
				Map(txs, func(txRow data.GetTxsRow) Node {
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
