package testutils

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var TestDBURL string = "postgres://test_user:test_pass@localhost:5433/test_db?sslmode=disable"

func ConnectToTestDB() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), TestDBURL)
}

func ApplyMigrations() error {
	m, err := migrate.New("file://../../../../migrations", TestDBURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}

	return nil
}
