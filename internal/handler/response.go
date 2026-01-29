// Утилиты для формирования HTTP-ответов.
package handler

import (
	"encoding/json"
	"net/http"
)

// JSONResponse отправляет JSON-ответ с указанным статусом.
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse отправляет ошибку в формате JSON.
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}

