package services

import (
	"log"

	"github.com/trieuvy/video-ranking/internal/models"
)

// processEvent processes a single event from the queue
type QueueServices struct {
	videoService *VideoService
	queue chan models.InteractionEvent
}

func NewQueueServices(videoService *VideoService, queue chan models.InteractionEvent) *QueueServices {
	return &QueueServices{videoService: videoService,queue: queue}
}

func (h *QueueServices) DequeueInteractionEvent(event models.InteractionEvent) {
	log.Printf("Processing event: %+v", event)
	if event.Type == models.Like {
		err := h.videoService.ChangeLikesAmount(event.VideoID, event.Step)
		if err != nil {
			log.Printf("Error changing likes amount: %v", err)
		}
	} else if event.Type == models.View {
		err := h.videoService.ChangeViewsAmount(event.VideoID, event.Step)
		if err != nil {
			log.Printf("Error changing views amount: %v", err)
		}
	} else if event.Type == models.Comment {
		err := h.videoService.ChangeCommentsAmount(event.VideoID, event.Step)
		if err != nil {
			log.Printf("Error changing comments amount: %v", err)
		}
	}
}

func (s *QueueServices) EnqueueInteractionEvent(event models.InteractionEvent) error {
	s.queue <- event
	return nil
}
