package service

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/product/model"
	"github.com/Just-maple/xmux/examples/webapp/internal/product/repository"
	"time"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.ProductResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	product := &model.Product{
		ID:          fmt.Sprintf("product-%d", time.Now().UnixNano()),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       100,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*model.ProductResponse, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, &model.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}

	return responses, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *model.UpdateProductRequest) (*model.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) error {
	return s.repo.Delete(ctx, req.ID)
}
