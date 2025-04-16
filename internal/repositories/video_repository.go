package repositories

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
	"gorm.io/gorm"
)

// VideoRepository handles database operations for videos
type VideoRepository struct {
	db *gorm.DB
}

// NewVideoRepository creates a new video repository
func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

// Create saves a new video to the database
func (r *VideoRepository) Create(video *models.Video) error {
	return r.db.Create(video).Error
}

// FindByID retrieves a video by ID
func (r *VideoRepository) FindByID(id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// FindByUser retrieves all videos by a user
func (r *VideoRepository) FindByUser(userID uuid.UUID) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Where("user_id = ?", userID).Find(&videos).Error
	return videos, err
}

// Update updates an existing video
func (r *VideoRepository) Update(video *models.Video) error {
	return r.db.Save(video).Error
}

// Delete removes a video from the database
func (r *VideoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Video{}, "id = ?", id).Error
}

// List retrieves all videos with pagination
func (r *VideoRepository) List(offset, limit int) ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Offset(offset).Limit(limit).Find(&videos).Error
	return videos, err
}

// UpdateScore updates the score of a video
func (r *VideoRepository) UpdateScore(id uuid.UUID, score float64) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).Update("score", score).Error
}

// IncrementLikes increments the like count for a video
func (r *VideoRepository) ChangeLikes(id uuid.UUID, step int) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).Update("likes", gorm.Expr("likes + ?", step)).Error
}

// IncrementViews increments the view count for a video
func (r *VideoRepository) ChangeViews(id uuid.UUID, step int) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).Update("views", gorm.Expr("views + ?", step)).Error
}

// IncrementComments increments the comment count for a video
func (r *VideoRepository) ChangeComments(id uuid.UUID, step int) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).Update("comments", gorm.Expr("comments + ?", step)).Error
}
