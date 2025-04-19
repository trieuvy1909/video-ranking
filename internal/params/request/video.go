package request

import (
	"github.com/google/uuid"
)

type Video struct {
	Title       string    `json:"title" validate:"required,min=3"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"create_by" validate:"required,uuid4"`
	Views       int64     `json:"views"`
	Likes       int64     `json:"likes"`
	Comments    int64     `json:"comments"`
	Score       float64   `json:"score"`
}
type VideoUpdate struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description"`
}