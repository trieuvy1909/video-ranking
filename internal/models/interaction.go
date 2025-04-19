package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InteractionType represents the type of interaction
type InteractionType string

const (
	Like    InteractionType = "like"
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

func (v *Interaction) BeforeUpdate(tx *gorm.DB) error {
	v.UpdatedAt = time.Now()
	return nil
}
func (v *Interaction) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return nil
}
