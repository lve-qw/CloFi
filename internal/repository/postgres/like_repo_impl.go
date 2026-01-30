package postgres

import (
	"context"

	"clofi/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLikeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) repository.LikeRepository {
	return &PostgresLikeRepository{db: db}
}

func (r *PostgresLikeRepository) AddLike(ctx context.Context, userID int64, productID string) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO likes (user_id, product_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, productID,
	)
	return err
}

func (r *PostgresLikeRepository) RemoveLike(ctx context.Context, userID int64, productID string) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM likes WHERE user_id = $1 AND product_id = $2",
		userID, productID,
	)
	return err
}

func (r *PostgresLikeRepository) IsLiked(ctx context.Context, userID int64, productID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND product_id = $2)",
		userID, productID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresLikeRepository) GetUserLikedProductIDs(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.Query(ctx,
		"SELECT product_id FROM likes WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
