package money

import (
	"github.com/keeth/money"
	"github.com/labstack/echo/v4"
)

func WithApp(app *money.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{
				Context: c,
				App:     app,
			}
			return next(cc)
		}
	}
}
