package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	money "github.com/keeth/money"
	data "github.com/keeth/money/data"
	"github.com/labstack/echo/v4"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

func main() {
	slog.SetDefault(defaultLogger)

	db, err := sql.Open("sqlite3", "file:money.db")
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}

	queries := data.New(db)

	money.InitGlobalApp(context.Background(), queries)

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
