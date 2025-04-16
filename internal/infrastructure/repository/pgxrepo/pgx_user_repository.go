package pgxrepo

import (
	"context"
	"errors"
	"rest-service/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository (db *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{db: db}
}

func (ur *PgxUserRepository) CreateUser (ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (id, refresh_token) VALUES ($1, $2 )`

	if _, err := ur.db.Exec(ctx, query, user.ID, user.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (ur *PgxUserRepository) GetUser (ctx context.Context, userID string) (*entities.User, error){
	var user entities.User
	query := `SELECT id, refresh_token FROM users WHERE id = $1`

	row := ur.db.QueryRow(ctx, query, userID)
	if err := row.Scan(&user.ID, &user.RefreshToken); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *PgxUserRepository) UpdateUser (ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET refresh_token = $2 WHERE id = $1`

	tag, err := ur.db.Exec(ctx, query, user.ID, user.RefreshToken)
	if tag.RowsAffected() == 0 {
		return errors.New("user doesn't exist")
	}

	return err
}