package repository

import (
	"context"
	"fmt"
	"github.com/Just-maple/xmux/examples/webapp/internal/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	users map[string]*model.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make(map[string]*model.User),
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}
