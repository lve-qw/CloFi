package model

type Like struct {
	UserID    int64  `db:"user_id"`
	ProductID string `db:"product_id"`
}
