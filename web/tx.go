package money

import (
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetTxs(c echo.Context) error {
	// app := money.GetGlobalApp()
	// txs, err := app.Queries.GetTxs(c.Request().Context(), data.GetTxsParams{
	// 	Ord:   "9999-99-99",
	// 	Limit: 10,
	// })
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, err.Error())
	// }
	page(PageProps{
		Title:       "Transactions",
		Description: "Transactions",
	},
		Div(Text("Hello, World!")),
	).Render(c.Response().Writer)
	return nil
}
