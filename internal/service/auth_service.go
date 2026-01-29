// Пакет service содержит бизнес-логику приложения.
package service

import (
	"context"
	"errors"

	"clofi/internal/model"
	"clofi/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("пользователь с таким логином уже существует")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")
)

// AuthService отвечает за регистрацию и аутентификацию.
type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register регистрирует нового пользователя.
// Возвращает ошибку, если логин уже занят.
func (s *AuthService) Register(ctx context.Context, req model.CreateUserRequest) error {
	// Проверяем, существует ли пользователь
	existing, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(ctx, user)
}

// Login проверяет учётные данные и возвращает пользователя.
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Сравниваем хеш
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}


