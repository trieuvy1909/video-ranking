package params

import (
	"time"
 	"gorm.io/gorm"
	"github.com/google/uuid"
)

// InteractionType represents the type of interaction
type InteractionType string

const (
	Like    InteractionType = "like"
	Dislike InteractionType = "dislike"
	View    InteractionType = "view"
	Comment InteractionType = "comment"
)

// Interaction represents a user interaction with a video
type Interaction struct {
	ID        uuid.UUID       `json:"id" gorm:"type:char(36);primary_key;"`
	UserID    uuid.UUID       `json:"user_id" gorm:"type:char(36);not null"`
	VideoID   uuid.UUID       `json:"video_id" gorm:"type:char(36);not null"`
	Type      InteractionType `json:"type" gorm:"size:20;not null"`
	Content   string          `json:"content" gorm:"type:text"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
func (v *Interaction) BeforeCreate(tx *gorm.DB) error {
    if v.ID == uuid.Nil {
        v.ID = uuid.New()
    }
    return nil
}
