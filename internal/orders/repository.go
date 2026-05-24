package orders

import (
	"context"

	"github.com/konkerama/go-grpc-api/internal/db"
)

type Repository struct {
	pool db.DBPool
}

func NewRepository(pool db.DBPool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	// 1. Data Persistence: Save the order to the database
	var createdOrder Order
	err := r.pool.QueryRow(ctx, `
		INSERT INTO orders (product_name, quantity)
		VALUES ($1, $2)
		RETURNING id, product_name, quantity
	`, order.ProductName, order.Quantity).Scan(&createdOrder.Id, &createdOrder.ProductName, &createdOrder.Quantity)
	if err != nil {
		return nil, err
	}

	// 2. Return the created order with its new ID
	return &createdOrder, nil
}
