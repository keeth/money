package money

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/terinjokes/sqlitestdb"
)

func TestNew(t *testing.T) {
	t.Parallel()
	conf := sqlitestdb.Config{Driver: "sqlite3"}

	migrator := sqlitestdb.NoopMigrator{}
	db := sqlitestdb.New(t, conf, migrator)

	var message string
	err := db.QueryRow("SELECT 'hellorld!'").Scan(&message)
	if err != nil {
		t.Fatalf("expected nil error: %+v\n", err)
	}

	if message != "hellord!" {
		t.Fatalf("expected message to be 'hellord!'")
	}
}
