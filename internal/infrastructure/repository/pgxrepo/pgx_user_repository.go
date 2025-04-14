package pgxrepo

import (
	"rest-service/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository (db *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{db: db}
}

func (ur *PgxUserRepository) CreateUser (user *entities.User) error {
	return nil
}