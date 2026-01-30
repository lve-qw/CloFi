package repository

import (
	"context"

	"clofi/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindAll(ctx context.Context, filters ProductFilters, page, limit int) ([]*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
	SearchByText(ctx context.Context, query string, page, limit int) ([]*model.Product, error)
}

type ProductFilters struct {
	Brand        *string // nil = без фильтра
	Availability *bool   // nil = без фильтра
	SortByPrice  string  // "asc", "desc" или пусто
}
