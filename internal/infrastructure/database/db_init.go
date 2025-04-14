package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDB() *pgxpool.Pool {
	dsn :="host=" + os.Getenv("POSTGRES_HOST") + 
	" user=" + os.Getenv("POSTGRES_USER") + 
	" password=" + os.Getenv("POSTGRES_PASSWORD") + 
	" dbname=" + os.Getenv("POSTGRES_DB") + 
	" port=" + os.Getenv("POSTGRES_PORT")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to open db connection; err: %v", err)
	}

	return pool
}