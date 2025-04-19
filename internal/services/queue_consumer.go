package services

import (
	"time"

	"github.com/trieuvy/video-ranking/internal/models"
)

// processEvent processes a single event from the queue
func QueueConsumer(queue chan models.InteractionEvent, QueueServices *QueueServices) {
	for {
		select {
		case event := <-queue:
			QueueServices.DequeueInteractionEvent(event)
		default:
			time.Sleep(1 * time.Second) // Sleep to prevent busy waiting
		}
	}
}