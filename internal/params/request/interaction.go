package request

import (
	"github.com/google/uuid"
	"github.com/trieuvy/video-ranking/internal/models"
)

// Interaction represents a user interaction with a video
type Interaction struct {
	UserID  uuid.UUID              `json:"user_id" validate:"required,uuid4"`
	VideoID uuid.UUID              `json:"video_id" validate:"required,uuid4"`
	Type    models.InteractionType `json:"type" validate:"required,oneof=like comment view"`
	Content string                 `json:"content" validate:"omitempty,max=1000"`
}
type InteractionUpdate struct {
	Content string `json:"content" validate:"omitempty,max=1000"`
}
