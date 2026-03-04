package repository

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/order/model"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type orderRepository struct {
	orders map[string]*model.Order
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	for _, order := range r.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	order, exists := r.orders[id]
	if !exists {
		return fmt.Errorf("order not found")
	}
	order.Status = status
	order.CreatedAt = time.Now()
	return nil
}
