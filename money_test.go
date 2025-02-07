package money

import (
	"context"
	"os"
	"testing"

	data "github.com/keeth/money/data"
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
	queries := data.New(db)
	app := NewApp(ctx, queries)
	ofxFile, err := os.Open("testdata/ofx/bank1.qfx")
	defer ofxFile.Close()
	assert.NoError(t, err)
	app.ImportOFX(ofxFile)
	acc, err := queries.GetAccs(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(acc))
	assert.Equal(t, "bank", acc[0].Name)
	assert.Equal(t, "000000001 001 bankAccount1234567890", acc[0].Xid)
	assert.Equal(t, "bank", acc[0].Kind)
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tx", acc[0].ID)
	var count int
	err = row.Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestImportSuite(t *testing.T) {
	suite.Run(t, new(ImportTestSuite))
}
