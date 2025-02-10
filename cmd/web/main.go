package main

import (
	"database/sql"
	"log/slog"
	"os"

	money "github.com/keeth/money"
	data "github.com/keeth/money/data"
	web "github.com/keeth/money/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

func main() {
	slog.SetDefault(defaultLogger)

	db, err := sql.Open("sqlite3", "file:money.db")
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}

	money.InitGlobalApp(data.New(db))

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	e := echo.New()
	e.Use(middleware.Static("static"))
	e.GET("/", web.GetIndex)
	e.GET("/tx", web.GetTxs)
	e.Logger.Fatal(e.Start(":" + portStr))
}
