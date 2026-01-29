package service

import (
	"context"
	"errors"

	"clofi/internal/repository"
)

var (
	ErrProductNotFound = errors.New("товар не найден")
	ErrUserNotFound    = errors.New("пользователь не найден")
)

// LikeService управляет лайками.
type LikeService struct {
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	likeRepo    repository.LikeRepository
}

func NewLikeService(
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	likeRepo repository.LikeRepository,
) *LikeService {
	return &LikeService{
		productRepo: productRepo,
		userRepo:    userRepo,
		likeRepo:    likeRepo,
	}
}

// ToggleLike добавляет или удаляет лайк у товара.
// Возвращает true, если лайк добавлен, false — если удалён.
func (s *LikeService) ToggleLike(ctx context.Context, userID int64, productID string) (bool, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrUserNotFound
	}

	// Проверяем существование товара
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return false, err
	}
	if product == nil {
		return false, ErrProductNotFound
	}

	// Проверяем, лайкал ли пользователь этот товар
	liked, err := s.likeRepo.IsLiked(ctx, userID, productID)
	if err != nil {
		return false, err
	}

	if liked {
		// Удаляем лайк
		err = s.likeRepo.RemoveLike(ctx, userID, productID)
		return false, err
	} else {
		// Добавляем лайк
		err = s.likeRepo.AddLike(ctx, userID, productID)
		return true, err
	}
}

// IsLiked проверяет, лайкал ли пользователь товар.
func (s *LikeService) IsLiked(ctx context.Context, userID int64, productID string) (bool, error) {
	return s.likeRepo.IsLiked(ctx, userID, productID)
}


