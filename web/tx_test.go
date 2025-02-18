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
	gm := golangmigrator.New("../db/migrations")
	db := sqlitestdb.New(t, sqlitestdb.Config{Driver: "sqlite3"}, gm)
	app := money.InitGlobalApp(db)
	ofxFile, err := os.Open("../testdata/ofx/bank1.qfx")
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
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(b.String()))
	t.NoError(err)
	t.Equal(1, doc.Find("html").Length())
	t.Equal(1, doc.Find("table").Length())
	rows := doc.Find("tr")
	t.Equal(3, rows.Length())
	lastRow := rows.Last()
	t.Equal("/tx?before=2025-02-02+4447", lastRow.AttrOr("hx-get", ""))
	t.Equal("afterend", lastRow.AttrOr("hx-swap", ""))
	t.Equal("revealed", lastRow.AttrOr("hx-trigger", ""))

	err, txs = GetTxs(ctx, model.GetTxsParams{
		Before: "2025-02-02 4447",
		Limit:  2,
	})
	t.NoError(err)
	var b2 strings.Builder
	t.NoError(txs.Render(&b2))
	// net/html requires a full html document
	doc2, err := goquery.NewDocumentFromReader(strings.NewReader("<html><head></head><body><table>" + b2.String() + "</table></body></html>"))
	t.NoError(err)
	t.Equal(1, doc2.Find("table").Length())
	rows2 := doc2.Find("tr")
	t.Equal(1, rows2.Length())
	lastRow2 := rows2.Last()
	t.Equal("/tx?before=2025-02-02+4446", lastRow2.AttrOr("hx-get", ""))
}

func TestWebTxTestSuite(t *testing.T) {
	suite.Run(t, new(WebTxTestSuite))
}
