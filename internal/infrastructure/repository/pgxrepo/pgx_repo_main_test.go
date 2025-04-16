package pgxrepo_test

import (
	"log"
	"os"
	testutils "rest-service/testutils/test_database"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var TestPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	TestPool, err = testutils.ConnectToTestDB()
	if err != nil {
		log.Fatalf("failed to open test db conenction; err: %v", err)
	}
	defer TestPool.Close()

	if err := testutils.ApplyMigrations(); err != nil {
		log.Fatalf("failed to apply migrations; err: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}
