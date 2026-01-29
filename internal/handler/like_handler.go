package handler

import (
	"net/http"

	"clofi/internal/service"
	"clofi/pkg/middleware"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// ToggleLike обрабатывает переключение лайка.
func (h *LikeHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(int64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "пользователь не авторизован")
		return
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		ErrorResponse(w, http.StatusBadRequest, "требуется product_id")
		return
	}

	liked, err := h.likeService.ToggleLike(r.Context(), userID, productID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrProductNotFound:
			ErrorResponse(w, http.StatusNotFound, err.Error())
		default:
			ErrorResponse(w, http.StatusInternalServerError, "ошибка обработки лайка")
		}
		return
	}

	status := "лайк добавлен"
	if !liked {
		status = "лайк удалён"
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": status})
}


