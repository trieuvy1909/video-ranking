package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/params/request"
	"github.com/trieuvy/video-ranking/internal/services"
)

// InteractionHandler handles HTTP requests for interactions
// @title Interaction API
// @description API for managing video interactions
type InteractionHandler struct {
	interactionService *services.InteractionService
	videoService       *services.VideoService
	userService       *services.UserService
}

// NewInteractionHandler creates a new interaction handler
func NewInteractionHandler(interactionService *services.InteractionService, videoService *services.VideoService, userService *services.UserService) *InteractionHandler {
	return &InteractionHandler{
		interactionService: interactionService,
		videoService:       videoService,
		userService:       userService,
	}
}

// CreateInteraction handles the creation of a new interaction
// @Summary Create a new interaction
// @Description Create a new interaction between a user and a video
// @Tags interactions
// @Accept json
// @Produce json
// @Param interaction body request.Interaction true "Interaction object"
// @Success 200 {object} models.Interaction
// @Failure 400 {string} string "Bad request"
// @Router /interactions [post]
func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction request.Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	err := validate.Struct(interaction)
	if err != nil {
		var sb strings.Builder
		for _, e := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on the '%s' rule\n", e.Field(), e.Tag()))
		}
		http.Error(w, sb.String(), http.StatusBadRequest)
		return
	}
	interactionModel := models.Interaction{
		UserID:  interaction.UserID,
		VideoID: interaction.VideoID,
		Type:    interaction.Type,
	}
	
	// Check if UserID exists
	if _, err := h.userService.GetUser(interactionModel.UserID); err!= nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Check if VideoID exists
	if _, err := h.videoService.GetVideo(interaction.VideoID); err!= nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	if err := h.interactionService.CreateInteraction(&interactionModel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if interaction.Type == models.Like {
		if err := h.videoService.ChangeLikesAmount(interaction.VideoID, 1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.View {
		if err := h.videoService.ChangeViewsAmount(interaction.VideoID, 1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.Comment {
		if err := h.videoService.ChangeCommentsAmount(interaction.VideoID, 1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// GetInteraction handles retrieving an interaction by ID
// @Summary Get an interaction by ID
// @Description Get details of a specific interaction
// @Tags interactions
// @Accept json
// @Produce json
// @Param id path string true "Interaction ID"
// @Success 200 {object} request.Interaction
// @Failure 400 {string} string "Invalid interaction ID"
// @Failure 404 {string} string "Interaction not found"
// @Router /interactions/{id} [get]
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
// @Summary Get all interactions between a user and a video
// @Description Get all interactions for a specific user and video combination
// @Tags interactions
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param videoID path string true "Video ID"
// @Success 200 {array} models.Interaction
// @Failure 400 {string} string "Invalid user or video ID"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{userID}/videos/{videoID}/interactions [get]
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
// @Summary Update an interaction
// @Description Update an existing interaction
// @Tags interactions
// @Accept json
// @Produce json
// @Param id path string true "Interaction ID"
// @Param interaction body request.Interaction true "Updated interaction object"
// @Success 200 {object} request.Interaction
// @Failure 400 {string} string "Invalid interaction ID or data"
// @Failure 500 {string} string "Internal server error"
// @Router /interactions/{id} [put]
func (h *InteractionHandler) UpdateInteraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid interaction ID", http.StatusBadRequest)
		return
	}
	var interaction request.InteractionUpdate
	if err = json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	err = validate.Struct(interaction)
	if err != nil {
		var sb strings.Builder
		for _, e := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on the '%s' rule\n", e.Field(), e.Tag()))
		}
		http.Error(w, sb.String(), http.StatusBadRequest)
		return
	}
	interactionModel, err := h.interactionService.GetInteraction(id)
	if err!= nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	interactionModel.Content =  interaction.Content
	if err := h.interactionService.UpdateInteraction(interactionModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// DeleteInteraction handles removing an interaction
// @Summary Delete an interaction
// @Description Delete an existing interaction
// @Tags interactions
// @Accept json
// @Produce json
// @Param id path string true "Interaction ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid interaction ID"
// @Failure 500 {string} string "Internal server error"
// @Router /interactions/{id} [delete]
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
		if err := h.videoService.ChangeLikesAmount(interaction.VideoID, -1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.View {
		if err := h.videoService.ChangeViewsAmount(interaction.VideoID, -1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if interaction.Type == models.Comment {
		if err := h.videoService.ChangeCommentsAmount(interaction.VideoID, -1); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListInteractions handles retrieving a list of interactions with pagination
// @Summary List all interactions
// @Description Get a paginated list of all interactions
// @Tags interactions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {array} models.Interaction
// @Failure 500 {string} string "Internal server error"
// @Router /interactions [get]
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
