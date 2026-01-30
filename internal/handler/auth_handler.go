package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"clofi/internal/model"
	"clofi/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
	jwtTTL      time.Duration
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string, jwtTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
		jwtTTL:      jwtTTL,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Register(r.Context(), req); err != nil {
		if err == service.ErrUserAlreadyExists {
			ErrorResponse(w, http.StatusConflict, err.Error())
		} else {
			ErrorResponse(w, http.StatusInternalServerError, "ошибка регистрации")
		}
		return
	}

	JSONResponse(w, http.StatusCreated, map[string]string{"message": "пользователь создан"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var cred struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	if err := h.validate.Struct(cred); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Login(r.Context(), cred.Username, cred.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			ErrorResponse(w, http.StatusUnauthorized, err.Error())
		} else {
			ErrorResponse(w, http.StatusInternalServerError, "ошибка входа")
		}
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(h.jwtTTL).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "ошибка генерации токена")
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"token": tokenStr})
}
