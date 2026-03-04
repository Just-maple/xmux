package service

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/order/model"
	orderRepo "github.com/Just-maple/xmux/examples/webapp/internal/order/repository"
	productRepo "github.com/Just-maple/xmux/examples/webapp/internal/product/repository"
	"time"
)

type OrderService struct {
	orderRepo   orderRepo.OrderRepository
	productRepo productRepo.ProductRepository
}

func NewOrderService(orderRepo orderRepo.OrderRepository, productRepo productRepo.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order items cannot be empty")
	}

	var total float64
	items := make([]model.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}

		itemPrice := product.Price * float64(item.Quantity)
		total += itemPrice

		items = append(items, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
		})

		product.Stock -= item.Quantity
		s.productRepo.Update(ctx, product)
	}

	order := &model.Order{
		ID:        fmt.Sprintf("order-%d", time.Now().UnixNano()),
		UserID:    req.UserID,
		Items:     items,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return &model.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *model.GetOrderRequest) (*model.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &model.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string) ([]*model.OrderResponse, error) {
	orders, err := s.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.OrderResponse, 0, len(orders))
	for _, o := range orders {
		responses = append(responses, &model.OrderResponse{
			ID:        o.ID,
			UserID:    o.UserID,
			Items:     o.Items,
			Total:     o.Total,
			Status:    o.Status,
			CreatedAt: o.CreatedAt,
		})
	}

	return responses, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	return s.orderRepo.UpdateStatus(ctx, id, status)
}
