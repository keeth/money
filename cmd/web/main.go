package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

func main() {
	slog.SetDefault(defaultLogger)

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + portStr))
}
