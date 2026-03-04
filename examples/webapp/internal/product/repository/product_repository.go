package repository

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/product/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id string) (*model.Product, error)
	GetAll(ctx context.Context) ([]*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
}

type productRepository struct {
	products map[string]*model.Product
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		products: make(map[string]*model.Product),
	}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	if _, exists := r.products[product.ID]; exists {
		return fmt.Errorf("product already exists")
	}
	r.products[product.ID] = product
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]*model.Product, error) {
	products := make([]*model.Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	if _, exists := r.products[product.ID]; !exists {
		return fmt.Errorf("product not found")
	}
	r.products[product.ID] = product
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.products[id]; !exists {
		return fmt.Errorf("product not found")
	}
	delete(r.products, id)
	return nil
}
