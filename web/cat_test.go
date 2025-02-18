package money

import (
	"context"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	money "github.com/keeth/money"
	model "github.com/keeth/money/model"
	sqlc "github.com/keeth/money/model/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"github.com/terinjokes/sqlitestdb"
	"github.com/terinjokes/sqlitestdb/migrators/golangmigrator"
)

type WebCatTestSuite struct {
	suite.Suite
	app *money.App
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *WebCatTestSuite) SetupTest() {
	t := suite.T()
	ctx := context.Background()
	gm := golangmigrator.New("../db/migrations")
	db := sqlitestdb.New(t, sqlitestdb.Config{Driver: "sqlite3"}, gm)
	suite.app = money.NewApp(db)
	for _, cat := range []string{"groceries", "rent", "utilities"} {
		_, err := suite.app.Model.Queries.CreateCat(ctx, sqlc.CreateCatParams{
			Name: cat,
			Kind: "expense",
		})
		suite.NoError(err)
	}
}

func (t *WebCatTestSuite) TestGetCats() {
	ctx := context.Background()
	err, cats := GetCats(ctx, t.app, model.GetCatsParams{
		After: "",
		Limit: 2,
	})
	t.NoError(err)
	var b strings.Builder
	t.NoError(cats.Render(&b))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(b.String()))
	t.NoError(err)
	t.Equal(1, doc.Find("html").Length())
	t.Equal(1, doc.Find("table").Length())
	rows := doc.Find("tr")
	t.Equal(3, rows.Length())
	lastRow := rows.Last()
	t.Equal("/cat?after=rent", lastRow.AttrOr("hx-get", ""))
	t.Equal("afterend", lastRow.AttrOr("hx-swap", ""))
	t.Equal("revealed", lastRow.AttrOr("hx-trigger", ""))

	err, cats = GetCats(ctx, t.app, model.GetCatsParams{
		After: "rent",
		Limit: 2,
	})
	t.NoError(err)
	var b2 strings.Builder
	t.NoError(cats.Render(&b2))
	// net/html requires a full html document
	doc2, err := goquery.NewDocumentFromReader(strings.NewReader("<html><head></head><body><table>" + b2.String() + "</table></body></html>"))
	t.NoError(err)
	t.Equal(1, doc2.Find("table").Length())
	rows2 := doc2.Find("tr")
	t.Equal(1, rows2.Length())
	lastRow2 := rows2.Last()
	t.Equal("/cat?after=utilities", lastRow2.AttrOr("hx-get", ""))
}

func TestWebCatTestSuite(t *testing.T) {
	suite.Run(t, new(WebCatTestSuite))
}
