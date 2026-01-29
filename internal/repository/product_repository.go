// Пакет repository содержит интерфейсы и реализации для доступа к данным.
package repository

import (
	"context"

	"clofi/internal/model"
)

// ProductRepository определяет методы для работы с товарами.
type ProductRepository interface {
	// Create добавляет новый товар в базу.
	Create(ctx context.Context, product *model.Product) error

	// FindAll возвращает список товаров с пагинацией и фильтрацией.
	FindAll(ctx context.Context, filters ProductFilters, page, limit int) ([]*model.Product, error)

	// FindByID ищет товар по ID.
	FindByID(ctx context.Context, id string) (*model.Product, error)

	// SearchByText выполняет текстовый поиск по name и description.
	SearchByText(ctx context.Context, query string, page, limit int) ([]*model.Product, error)
}

// ProductFilters содержит параметры фильтрации и сортировки.
type ProductFilters struct {
	Brand        *string // nil = без фильтра
	Availability *bool   // nil = без фильтра
	SortByPrice  string  // "asc", "desc" или пусто
}

