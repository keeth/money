package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	core "github.com/keeth/money/core"
	sqlc "github.com/keeth/money/model/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
var app *core.App

func main() {
	slog.SetDefault(defaultLogger)

	db, err := sql.Open("sqlite3", "file:money.db")
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}

	app = core.NewApp(db)

	var rootCmd = &cobra.Command{
		Use:   "money",
		Short: "Money management CLI",
	}

	var importCmd = &cobra.Command{
		Use:   "import [file]",
		Short: "Import transactions from OFX file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			result, err := app.ImportOFX(context.Background(), file)
			if err != nil {
				return err
			}

			slog.Info("import completed",
				"transactions_created", result.TxCreated,
				"transactions_updated", result.TxUpdated,
				"accounts_created", result.AccCreated)
			return nil
		},
	}

	var createCatCmd = &cobra.Command{
		Use:   "create-cat [name] [kind]",
		Short: "Create a new category",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := app.Model.Queries.CreateCat(context.Background(), sqlc.CreateCatParams{
				Name: args[0],
				Kind: args[1],
			})
			if err != nil {
				return err
			}
			slog.Info("category created", "name", args[0], "kind", args[1], "id", id)
			return nil
		},
	}

	var createRuleCmd = &cobra.Command{
		Use:   "create-rule [test]",
		Short: "Create a new rule",
		Long:  "Create a new rule with a test expression and optional category, amount, description, and date patterns. At least one optional parameter must be specified.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			category, _ := cmd.Flags().GetString("cat")
			amount, _ := cmd.Flags().GetString("amount")
			desc, _ := cmd.Flags().GetString("desc")
			date, _ := cmd.Flags().GetString("date")

			// Verify at least one optional parameter is provided
			if category == "" && amount == "" && desc == "" && date == "" {
				return fmt.Errorf("at least one of --cat, --amount, --desc, or --date must be specified")
			}

			var catID sql.NullInt64
			if category != "" {
				cat, err := app.Model.Queries.GetCatByName(context.Background(), category)
				if err != nil {
					return err
				}
				if cat.ID == 0 {
					return fmt.Errorf("category not found: %s", category)
				}
				catID = sql.NullInt64{Int64: cat.ID, Valid: true}
			}

			id, err := app.Model.CreateRule(context.Background(), sqlc.CreateRuleParams{
				TestExpr:   args[0],
				CatID:      catID,
				AmountExpr: sql.NullString{String: amount, Valid: amount != ""},
				DescExpr:   sql.NullString{String: desc, Valid: desc != ""},
				DateExpr:   sql.NullString{String: date, Valid: date != ""},
			})
			if err != nil {
				return err
			}
			slog.Info("rule created", "id", id)
			return nil
		},
	}

	createRuleCmd.Flags().String("cat", "", "Category ID to assign")
	createRuleCmd.Flags().String("amount", "", "Amount pattern")
	createRuleCmd.Flags().String("desc", "", "Description pattern")
	createRuleCmd.Flags().String("date", "", "Date pattern")

	rootCmd.AddCommand(importCmd, createCatCmd, createRuleCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}
