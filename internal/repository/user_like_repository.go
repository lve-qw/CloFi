package repository

import (
	"context"

	"clofi/internal/model"
)

// UserRepository — интерфейс для работы с пользователями.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

// LikeRepository — интерфейс для работы с лайками.
type LikeRepository interface {
	AddLike(ctx context.Context, userID int64, productID string) error
	RemoveLike(ctx context.Context, userID int64, productID string) error
	IsLiked(ctx context.Context, userID int64, productID string) (bool, error)
	GetUserLikedProductIDs(ctx context.Context, userID int64) ([]string, error)
}


