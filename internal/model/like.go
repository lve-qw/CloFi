package model

// Like связывает пользователя и товар (лайк).
// Хранится в PostgreSQL.
type Like struct {
	UserID    int64  `db:"user_id"`
	ProductID string `db:"product_id"`
}

