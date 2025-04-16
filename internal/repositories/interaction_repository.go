package repositories

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
	"gorm.io/gorm"
)

// InteractionRepository handles database operations for interactions
type InteractionRepository struct {
	db *gorm.DB
}

// NewInteractionRepository creates a new interaction repository
func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

// Create saves a new interaction to the database
func (r *InteractionRepository) Create(interaction *models.Interaction) error {
	return r.db.Create(interaction).Error
}

// FindByID retrieves an interaction by ID
func (r *InteractionRepository) FindByID(id uuid.UUID) (*models.Interaction, error) {
	var interaction models.Interaction
	err := r.db.First(&interaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &interaction, nil
}

// FindByUserAndVideo retrieves all interactions between a user and a video
func (r *InteractionRepository) FindByUserAndVideo(userID, videoID uuid.UUID) ([]models.Interaction, error) {
	var interactions []models.Interaction
	err := r.db.Where("user_id = ? AND video_id = ?", userID, videoID).Find(&interactions).Error
	return interactions, err
}

// Update updates an existing interaction
func (r *InteractionRepository) Update(interaction *models.Interaction) error {
	return r.db.Save(interaction).Error
}

// Delete removes an interaction from the database
func (r *InteractionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Interaction{}, "id = ?", id).Error
}

// List retrieves all interactions with pagination
func (r *InteractionRepository) List(offset, limit int) ([]models.Interaction, error) {
	var interactions []models.Interaction
	err := r.db.Offset(offset).Limit(limit).Find(&interactions).Error
	return interactions, err
}

// FindByVideo retrieves all interactions for a specific video
func (r *InteractionRepository) FindByVideo(videoID uuid.UUID) ([]models.Interaction, error) {
	var interactions []models.Interaction
	err := r.db.Where("video_id = ?", videoID).Find(&interactions).Error
	return interactions, err
}

// FindByUser retrieves all interactions for a specific user
func (r *InteractionRepository) FindByUser(userID uuid.UUID) ([]models.Interaction, error) {
	var interactions []models.Interaction
	err := r.db.Where("user_id = ?", userID).Find(&interactions).Error
	return interactions, err
}
