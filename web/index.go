package money

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getIndex(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
