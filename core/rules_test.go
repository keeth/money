package money

import (
	"context"
	"database/sql"
	"os"
	"testing"

	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/terinjokes/sqlitestdb"
	"github.com/terinjokes/sqlitestdb/migrators/golangmigrator"
)

type RulesTestSuite struct {
	suite.Suite
	app *App
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *RulesTestSuite) SetupTest() {
	gm := golangmigrator.New("../db/migrations")
	db := sqlitestdb.New(suite.T(), sqlitestdb.Config{Driver: "sqlite3"}, gm)
	suite.app = NewApp(db)
	ctx := context.Background()
	for _, cat := range []string{"groceries", "rent", "utilities"} {
		_, err := suite.app.Model.Queries.CreateCat(ctx, sqlc.CreateCatParams{
			Name: cat,
			Kind: "expense",
		})
		suite.NoError(err)
	}
}

func (suite *RulesTestSuite) TestApplyRules() {
	ctx := context.Background()
	app := suite.app
	groceriesCat, err := app.Model.Queries.GetCatByName(ctx, "groceries")
	assert.NoError(suite.T(), err)
	app.Model.CreateRule(ctx, model.CreateRuleParams{
		TestExpr: "Desc matches 'grocery'",
		CatID:    sql.NullInt64{Int64: groceriesCat.ID, Valid: true},
	})
	ofxFile, err := os.Open("../testdata/ofx/cc1.qfx")
	defer ofxFile.Close()
	assert.NoError(suite.T(), err)
	_, err = app.ImportOFX(ctx, ofxFile)
	assert.NoError(suite.T(), err)
	tx, err := app.Model.Queries.GetTxByAccAndXid(ctx, sqlc.GetTxByAccAndXidParams{
		AccID: 1,
		Xid:   "202503104597079",
	})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), sql.NullInt64{Int64: groceriesCat.ID, Valid: true}, tx.CatID)
}

func TestRulesTestSuite(t *testing.T) {
	suite.Run(t, new(RulesTestSuite))
}
