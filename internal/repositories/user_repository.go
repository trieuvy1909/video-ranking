package repositories

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create saves a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// List retrieves all users with pagination
func (r *UserRepository) List(page, pageSize int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, err
}
