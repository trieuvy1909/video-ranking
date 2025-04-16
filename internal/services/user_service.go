package services

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/repositories"
)

// UserService handles business logic for users
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(page, pageSize int) ([]models.User, error) {
	return s.userRepo.List(page, pageSize)
}
