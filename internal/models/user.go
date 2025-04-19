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
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate đảm bảo ID được tạo trước khi lưu vào database
func (v *User) BeforeUpdate(tx *gorm.DB) error {
	v.UpdatedAt = time.Now()
	return nil
}
func (v *User) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return nil
}
