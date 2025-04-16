package services

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/repositories"
)

type InteractionService struct {
	repo        *repositories.InteractionRepository
}

// NewInteractionService creates a new interaction service
func NewInteractionService(repo *repositories.InteractionRepository) *InteractionService {
	return &InteractionService{repo: repo}
}

// CreateInteraction creates a new interaction and processes it
func (s *InteractionService) CreateInteraction(interaction *models.Interaction) error {
	return s.repo.Create(interaction)
}

// GetInteraction retrieves an interaction by ID
func (s *InteractionService) GetInteraction(id uuid.UUID) (*models.Interaction, error) {
	return s.repo.FindByID(id)
}

// GetUserVideoInteractions retrieves all interactions between a user and a video
func (s *InteractionService) GetUserVideoInteractions(userID, videoID uuid.UUID) ([]models.Interaction, error) {
	return s.repo.FindByUserAndVideo(userID, videoID)
}

// UpdateInteraction updates an existing interaction
func (s *InteractionService) UpdateInteraction(interaction *models.Interaction) error {
	return s.repo.Update(interaction)
}

// DeleteInteraction removes an interaction
func (s *InteractionService) DeleteInteraction(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// ListInteractions retrieves a list of interactions with pagination
func (s *InteractionService) ListInteractions(page, pageSize int) ([]models.Interaction, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
