package money

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	money "github.com/keeth/money"
	model "github.com/keeth/money/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/terinjokes/sqlitestdb"
	"github.com/terinjokes/sqlitestdb/migrators/golangmigrator"
)

type WebTxTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *WebTxTestSuite) SetupTest() {
	t := suite.T()
	ctx := context.Background()
	gm := golangmigrator.New("db/migrations")
	db := sqlitestdb.New(t, sqlitestdb.Config{Driver: "sqlite3"}, gm)
	app := money.InitGlobalApp(db)
	ofxFile, err := os.Open("testdata/ofx/bank1.qfx")
	defer ofxFile.Close()
	assert.NoError(t, err)
	app.ImportOFX(ctx, ofxFile)
}

func (t *WebTxTestSuite) TestGetTxs() {
	ctx := context.Background()
	err, txs := GetTxs(ctx, model.GetTxsParams{
		Before: "",
		Limit:  2,
	})
	t.NoError(err)
	var b strings.Builder
	t.NoError(txs.Render(&b))
	println(b.String())
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(b.String()))
	t.NoError(err)
	t.Equal(1, doc.Find("html").Length())
	t.Equal(1, doc.Find("table").Length())
	rows := doc.Find("tr")
	t.Equal(3, rows.Length())
}

func TestWebTxTestSuite(t *testing.T) {
	suite.Run(t, new(WebTxTestSuite))
}
