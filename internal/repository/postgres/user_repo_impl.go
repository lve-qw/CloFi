package postgres

import (
	"context"

	"clofi/internal/model"
	"clofi/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO users (username, name, password) VALUES ($1, $2, $3)",
		user.Username, user.Name, user.Password,
	)
	return err
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, username, name, password FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Name, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		"SELECT id, username, name, password FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Name, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}


