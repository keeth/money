package money

import (
	"context"
	"log/slog"

	core "github.com/keeth/money/core"
	model "github.com/keeth/money/model"
	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GetCatRows(cats []model.Cat) Node {
	return Group{Map(cats, func(cat model.Cat) Node {
		return Tr(
			If(cat.ID == cats[len(cats)-1].ID,
				nextPageNode("cat", cat.Name),
			),
			Td(Text(cat.Name)),
			Td(Text(cat.Kind)),
		)
	})}
}

func GetCats(ctx context.Context, app *core.App, params model.GetCatsParams) (error, Node) {
	cats, err := app.Model.GetCats(ctx, params)
	if err != nil {
		slog.Error("failed to get cats", "error", err)
		return err, nil
	}
	rowNodes := GetCatRows(cats)
	if params.After == "" {
		return tablePage(PageProps{
			Title:       "Categories",
			Description: "Categories",
		},
			[]string{"Name", "Kind"},
			rowNodes,
		)
	}
	return nil, rowNodes
}

func GetCatsEndpoint(c echo.Context) error {
	cc := c.(*Context)
	after, limit := paginationParams(c)
	err, cats := GetCats(cc.Request().Context(), cc.App, model.GetCatsParams{
		After: after,
		Limit: limit,
	})
	return render(c, err, cats)
}
