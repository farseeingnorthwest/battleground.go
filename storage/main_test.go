package storage_test

import (
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db       *sqlx.DB
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		db = sqlx.MustConnect("postgres", url)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)

		var err error
		if fixtures, err = testfixtures.NewFolder(db.DB, &testfixtures.PostgreSQL{}, "fixtures"); err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}

func loadFixtures(t *testing.T) {
	if fixtures == nil {
		t.Skip("DATABASE_URL is not set")
	}
	if err := fixtures.Load(); err != nil {
		t.Fatal(err)
	}
}
