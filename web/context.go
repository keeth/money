package money

import (
	core "github.com/keeth/money/core"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	App *core.App
}
