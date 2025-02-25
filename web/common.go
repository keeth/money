package money

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
}

func page(props PageProps, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href("/css/cu.css")),
			Script(Src("https://unpkg.com/htmx.org")),
		},
		Body: []Node{
			Group(children),
		},
	})
}

var maxLimit = int64(100)

func paginationParams(c echo.Context) (string, int64) {
	after := c.QueryParam("after")
	limitStr := c.QueryParam("limit")
	limit := maxLimit
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return "", 0
		}
		if limit < 1 || limit > maxLimit {
			limit = maxLimit
		}
	}
	return after, limit
}

func paginationParamsInt(c echo.Context) (int64, int64) {
	after, limit := paginationParams(c)
	afterInt, err := strconv.ParseInt(after, 10, 64)
	if err != nil {
		return 0, limit
	}
	return afterInt, limit
}

func render(c echo.Context, err error, node Node) error {
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return node.Render(c.Response().Writer)
}

func tablePage(props PageProps, cols []string, rows Node) (error, Node) {
	return nil, page(props,
		Table(
			THead(
				Tr(
					Map(cols, func(col string) Node {
						return Th(Text(col))
					}),
				),
			),
			TBody(
				rows,
			),
		),
	)
}

func nextPageNode(path string, after string) Node {
	return Group{
		hx.Get(fmt.Sprintf("/%s?after=%s", path, url.QueryEscape(after))),
		hx.Trigger("revealed"),
		hx.Swap("afterend"),
	}
}
