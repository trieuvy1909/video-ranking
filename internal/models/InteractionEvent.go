package models

import (
	"github.com/google/uuid"
)
type InteractionEvent struct {
	VideoID   uuid.UUID 
	Type InteractionType
	Step int
}