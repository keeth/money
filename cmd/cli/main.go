package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/keeth/money"
	data "github.com/keeth/money/data"
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

	app := money.InitGlobalApp(data.New(db))

	if len(os.Args) < 2 {
		slog.Error("no command specified")
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "import" {
		if len(os.Args) < 3 {
			slog.Error("import command requires a file path")
			os.Exit(1)
		}

		filePath := os.Args[2]
		file, err := os.Open(filePath)
		if err != nil {
			slog.Error("failed to open import file", "err", err)
			os.Exit(1)
		}
		defer file.Close()

		result, err := app.ImportOFX(context.Background(), file)
		if err != nil {
			slog.Error("import failed", "err", err)
			os.Exit(1)
		}

		slog.Info("import completed",
			"transactions_created", result.TxCreated,
			"transactions_updated", result.TxUpdated,
			"accounts_created", result.AccCreated)
	} else {
		slog.Error("unknown command", "cmd", cmd)
		os.Exit(1)
	}
}
