package handler

import (
	"net/http"
	"strconv"
	"strings"

	"clofi/internal/repository"
	"clofi/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var filters repository.ProductFilters
	if brand := r.URL.Query().Get("brand"); brand != "" {
		filters.Brand = &brand
	}
	if availStr := r.URL.Query().Get("availability"); availStr != "" {
		avail := availStr == "true"
		filters.Availability = &avail
	}
	if sort := r.URL.Query().Get("sort_price"); sort != "" {
		filters.SortByPrice = sort
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	products, err := h.productService.GetProducts(r.Context(), filters, query, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "ошибка получения товаров")
		return
	}

	JSONResponse(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		ErrorResponse(w, http.StatusBadRequest, "требуется параметр id")
		return
	}

	if len(id) != 24 {
		ErrorResponse(w, http.StatusBadRequest, "некоректный id")
	}

	product, err := h.productService.GetProductByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "ошибка получения товара")
		return
	}
	if product == nil {
		ErrorResponse(w, http.StatusNotFound, "товар не найден")
		return
	}

	JSONResponse(w, http.StatusOK, product)
}
