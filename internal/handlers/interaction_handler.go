package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/services"
)

// InteractionHandler handles HTTP requests for interactions
type InteractionHandler struct {
	interactionService *services.InteractionService;
	videoService       *services.VideoService;
}

// NewInteractionHandler creates a new interaction handler
func NewInteractionHandler(interactionService *services.InteractionService, videoService *services.VideoService) *InteractionHandler {
	return &InteractionHandler{
		interactionService: interactionService,
		videoService:       videoService,
	}
}

// CreateInteraction handles the creation of a new interaction
func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction models.Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.interactionService.CreateInteraction(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if interaction.Type == models.Like {
		if err := h.videoService.ChangeLikesAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.View {
		if err := h.videoService.ChangeViewsAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.Comment {
		if err := h.videoService.ChangeCommentsAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// GetInteraction handles retrieving an interaction by ID
func (h *InteractionHandler) GetInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid interaction ID", http.StatusBadRequest)
		return
	}

	interaction, err := h.interactionService.GetInteraction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// GetUserVideoInteractions handles retrieving all interactions between a user and a video
func (h *InteractionHandler) GetUserVideoInteractions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	videoID, err := uuid.Parse(vars["videoID"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	interactions, err := h.interactionService.GetUserVideoInteractions(userID, videoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interactions)
}

// UpdateInteraction handles updating an existing interaction
func (h *InteractionHandler) UpdateInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid interaction ID", http.StatusBadRequest)
		return
	}

	var interaction models.Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	interaction.ID = id
	if err := h.interactionService.UpdateInteraction(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// DeleteInteraction handles removing an interaction
func (h *InteractionHandler) DeleteInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid interaction ID", http.StatusBadRequest)
		return
	}
	interaction, err := h.interactionService.GetInteraction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := h.interactionService.DeleteInteraction(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if interaction.Type == models.Like {
		if err := h.videoService.ChangeLikesAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.View {
		if err := h.videoService.ChangeViewsAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.Comment {
		if err := h.videoService.ChangeCommentsAmount(interaction.VideoID,1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListInteractions handles retrieving a list of interactions with pagination
func (h *InteractionHandler) ListInteractions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	interactions, err := h.interactionService.ListInteractions(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interactions)
}

// RegisterRoutes registers the interaction routes
func (h *InteractionHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/interactions", h.CreateInteraction).Methods("POST")
	r.HandleFunc("/interactions/{id}", h.GetInteraction).Methods("GET")
	r.HandleFunc("/interactions/{id}", h.UpdateInteraction).Methods("PUT")
	r.HandleFunc("/interactions/{id}", h.DeleteInteraction).Methods("DELETE")
	r.HandleFunc("/interactions", h.ListInteractions).Methods("GET")
	r.HandleFunc("/users/{userID}/videos/{videoID}/interactions", h.GetUserVideoInteractions).Methods("GET")
}
