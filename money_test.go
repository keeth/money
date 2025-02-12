package money

import (
	"context"
	"os"
	"testing"

	sqlc "github.com/keeth/money/model/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/terinjokes/sqlitestdb"
	"github.com/terinjokes/sqlitestdb/migrators/golangmigrator"
)

type ImportTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ImportTestSuite) SetupTest() {
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ImportTestSuite) TestImportOFX() {
	t := suite.T()
	ctx := context.Background()
	gm := golangmigrator.New("db/migrations")
	db := sqlitestdb.New(t, sqlitestdb.Config{Driver: "sqlite3"}, gm)
	app := InitGlobalApp(db)
	ofxFile, err := os.Open("testdata/ofx/bank1.qfx")
	defer ofxFile.Close()
	assert.NoError(t, err)
	app.ImportOFX(ctx, ofxFile)
	acc, err := app.Model.Queries.GetAccs(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(acc))
	assert.Equal(t, "bank", acc[0].Name)
	assert.Equal(t, "000000001 001 bankAccount1234567890", acc[0].Xid)
	assert.Equal(t, "bank", acc[0].Kind)
	txs, err := app.Model.Queries.GetTxs(ctx, sqlc.GetTxsParams{
		Ord:   "9999-99-99",
		Limit: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(txs))
	assert.Equal(t, "2025-02-02", txs[0].Tx.Date)
	assert.Equal(t, "emt transfer - credit payer: fozzie bear", txs[0].Tx.Desc)
	assert.Equal(t, 100.0, txs[0].Tx.Amount)
	assert.Equal(t, acc[0].ID, txs[0].Tx.AccID)
	assert.Equal(t, acc[0].ID, txs[0].Acc.ID)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestImportSuite(t *testing.T) {
	suite.Run(t, new(ImportTestSuite))
}
