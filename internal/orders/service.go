package orders

import (
	"context"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
}

type Service struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	// 2. Data Persistence: Save the order to the database
	createdOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// 3. Return the created order with its new ID
	return createdOrder, nil
}
