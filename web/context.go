package money

import (
	"github.com/keeth/money"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	App *money.App
}
