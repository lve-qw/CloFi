package service

import (
	"context"

	"clofi/internal/model"
	"clofi/internal/repository"
)

// ProductService управляет товарами: поиск, фильтрация, пагинация.
type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// GetProducts возвращает список товаров с фильтрацией и пагинацией.
func (s *ProductService) GetProducts(
	ctx context.Context,
	filters repository.ProductFilters,
	query string,
	page, limit int,
) ([]*model.Product, error) {
	if query != "" {
		return s.productRepo.SearchByText(ctx, query, page, limit)
	}
	return s.productRepo.FindAll(ctx, filters, page, limit)
}

// GetProductByID возвращает товар по ID.
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}


