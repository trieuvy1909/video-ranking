package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// Video represents a video entity in the system
type Video struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	URL         string    `json:"url" gorm:"size:255;not null"`
	CreatedBy   uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	Views       int64     `json:"views" gorm:"default:0"`
	Likes       int64     `json:"likes" gorm:"default:0"`
	Comments    int64     `json:"comments" gorm:"default:0"`
	Score       float64   `json:"score" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
func (v *Video) BeforeCreate(tx *gorm.DB) error {
    if v.ID == uuid.Nil {
        v.ID = uuid.New()
    }
    return nil
}
