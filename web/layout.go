package money

import (
	"html/template"
	"io"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	slog.Debug("Render", "name", name, "data", data)
	return t.templates.ExecuteTemplate(w, name, data)
}
