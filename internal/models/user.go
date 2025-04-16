package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User represents a user in the system
type User struct {
    ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
    Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
    Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"-"`
    AvatarURL string    `gorm:"size:255" json:"avatar_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate đảm bảo ID được tạo trước khi lưu vào database
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}
