package service

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/user/model"
	"github.com/Just-maple/xmux/examples/webapp/internal/user/repository"
	"time"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user := &model.User{
		ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *model.DeleteUserRequest) error {
	return s.repo.Delete(ctx, req.ID)
}
