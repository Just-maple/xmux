package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/Just-maple/xmux/examples/gin-business/internal/models"
	"github.com/Just-maple/xmux/examples/gin-business/internal/repository"
	"github.com/Just-maple/xmux/examples/gin-business/internal/types"
	"github.com/google/uuid"
)

// UserService defines the business logic interface for user operations
type UserService interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.UserResponse, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, userID uuid.UUID) (*types.UserResponse, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, userID uuid.UUID, req *types.UpdateUserRequest) (*types.UserResponse, error)

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// ListUsers lists users with pagination
	ListUsers(ctx context.Context, req *types.ListUsersRequest) (*types.ListUsersResponse, error)

	// ChangePassword changes user's password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *types.ChangePasswordRequest) error

	// Login authenticates a user and returns a token
	Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error)

	// GetProfile returns the current user's profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*types.UserResponse, error)
}

// userServiceImpl is the implementation of UserService
type userServiceImpl struct {
	userRepo repository.UserRepository
	// In a real application, you'd also have:
	// - authService for token generation
	// - emailService for sending emails
	// - logger for logging
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.UserResponse, error) {
	// Validate role
	role := models.UserRole(req.Role)
	if role != models.RoleAdmin && role != models.RoleUser && role != models.RoleViewer {
		return nil, errors.New("invalid role")
	}

	// Create user model
	user, err := models.NewUser(
		req.Username,
		req.Email,
		req.Password,
		req.FullName,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return types.FromUser(user), nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, userID uuid.UUID) (*types.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return types.FromUser(user), nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, userID uuid.UUID, req *types.UpdateUserRequest) (*types.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update user profile
	user.UpdateProfile(req.FullName, req.Email)

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return types.FromUser(user), nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *userServiceImpl) ListUsers(ctx context.Context, req *types.ListUsersRequest) (*types.ListUsersResponse, error) {
	// Get users with pagination
	users, err := s.userRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &types.ListUsersResponse{
		Users:   types.FromUsers(users),
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+req.Limit < total,
	}, nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, req *types.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Verify old password
	if !user.CheckPassword(req.OldPassword) {
		return repository.ErrInvalidPassword
	}

	// Update password
	if err := user.ChangePassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// Also try email
		user, err = s.userRepo.FindByEmail(ctx, req.Username)
		if err != nil {
			return nil, repository.ErrUserNotFound
		}
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, repository.ErrInvalidPassword
	}

	// In a real application, you would generate a JWT token here
	token := "mock-jwt-token-for-" + user.ID.String()

	return &types.LoginResponse{
		Token: token,
		User:  *types.FromUser(user),
	}, nil
}

func (s *userServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*types.UserResponse, error) {
	return s.GetUser(ctx, userID)
}
