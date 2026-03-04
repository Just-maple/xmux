package business

import (
	"context"
	"fmt"
	"time"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserRequest struct {
	ID string `json:"id"`
}

type UpdateUserRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserService struct {
	users map[string]*UserResponse
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*UserResponse),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	id := fmt.Sprintf("user-%d", time.Now().UnixNano())

	user := &UserResponse{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: time.Now(),
	}

	s.users[id] = user

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	user, exists := s.users[req.ID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", req.ID)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	user, exists := s.users[req.ID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", req.ID)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	_, exists := s.users[req.ID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", req.ID)
	}

	delete(s.users, req.ID)

	return &DeleteUserResponse{
		Success: true,
		Message: fmt.Sprintf("user %s deleted", req.ID),
	}, nil
}

type ListUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListUsersResponse struct {
	Users []*UserResponse `json:"users"`
	Total int             `json:"total"`
}

func (s *UserService) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	users := make([]*UserResponse, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	total := len(users)
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	start := req.Offset
	if start < 0 {
		start = 0
	}
	end := start + limit
	if end > total {
		end = total
	}

	var pageUsers []*UserResponse
	if start < total {
		pageUsers = users[start:end]
	} else {
		pageUsers = []*UserResponse{}
	}

	return &ListUsersResponse{
		Users: pageUsers,
		Total: total,
	}, nil
}
