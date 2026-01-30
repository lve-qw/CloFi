package service

import (
	"context"

	"clofi/internal/model"
	"clofi/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

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

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}
