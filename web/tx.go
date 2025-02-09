package money

import (
	"context"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

func getTxs(c echo.Context) error {
	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	ctx := context.Background()
	txs, err := a.queries.GetTxs(a.ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = tmp.ExecuteTemplate(c.Response().Writer, "index.html", txs)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return nil
}
