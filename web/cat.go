package money

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/keeth/money"
	model "github.com/keeth/money/model"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func GetCatRows(cats []model.Cat) Node {
	return Group{Map(cats, func(cat model.Cat) Node {
		return Tr(
			If(cat.ID == cats[len(cats)-1].ID,
				Group{
					hx.Get(fmt.Sprintf("/cat?after=%s", cat.Name)),
					hx.Trigger("revealed"),
					hx.Swap("afterend"),
				},
			),
			Td(Text(cat.Name)),
			Td(Text(cat.Kind)),
		)
	})}
}

func GetCats(ctx context.Context, app *money.App, params model.GetCatsParams) (error, Node) {
	cats, err := app.Model.GetCats(ctx, params)
	if err != nil {
		slog.Error("failed to get cats", "error", err)
		return err, nil
	}
	rowNodes := GetCatRows(cats)
	if params.After == "" {
		return nil, page(PageProps{
			Title:       "Categories",
			Description: "Categories",
		},
			Table(
				THead(
					Tr(
						Th(Text("Name")),
						Th(Text("Kind")),
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

func GetCatsEndpoint(c echo.Context) error {
	cc := c.(*Context)
	after := c.QueryParam("after")
	limitStr := c.QueryParam("limit")
	limit := maxLimit
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid limit")
		}
		if limit < 1 || limit > maxLimit {
			limit = maxLimit
		}
	}
	err, cats := GetCats(cc.Request().Context(), cc.App, model.GetCatsParams{
		After: after,
		Limit: limit,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return cats.Render(c.Response().Writer)
}
