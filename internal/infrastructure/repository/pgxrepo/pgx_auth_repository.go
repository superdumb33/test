package pgxrepo

import (
	"context"
	"rest-service/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxAuthRepository struct {
	db *pgxpool.Pool
}

func NewPgxAuthRepository (db *pgxpool.Pool) *PgxAuthRepository {
	return &PgxAuthRepository{db: db}
}

func (ar *PgxAuthRepository) CreateUser (ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (id, refresh_token) VALUES ($1, $2 )`

	if _, err := ar.db.Exec(ctx, query, user.ID, user.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (ar *PgxAuthRepository) GetUser (ctx context.Context, userID string) (*entities.User, error){
	var user entities.User
	query := `SELECT id, refresh_token FROM users WHERE id = $1`

	row := ar.db.QueryRow(ctx, query, userID)
	if err := row.Scan(&user.ID, &user.RefreshToken); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *PgxAuthRepository) UpdateUser (ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET refresh_token = $2 WHERE id = $1`

	tag, err := ar.db.Exec(ctx, query, user.ID, user.RefreshToken)
	if tag.RowsAffected() == 0 {
		return entities.NewAppErr(550, "no such user")
	}

	return err
}