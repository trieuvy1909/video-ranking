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

// VideoHandler handles HTTP requests for videos
// @title Video API
// @description API for managing videos
type VideoHandler struct {
	videoService *services.VideoService
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

// CreateVideo handles the creation of a new video
// @Summary Create a new video
// @Description Create a new video in the system
// @Tags videos
// @Accept json
// @Produce json
// @Param video body models.Video true "Video object"
// @Success 200 {object} models.Video
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /videos [post]
func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var video models.Video
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.videoService.CreateVideo(&video); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

// GetVideo handles retrieving a video by ID
// @Summary Get a video by ID
// @Description Get details of a specific video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} models.Video
// @Failure 400 {string} string "Invalid video ID"
// @Failure 404 {string} string "Video not found"
// @Router /videos/{id} [get]
func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	video, err := h.videoService.GetVideo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

// GetUserVideos handles retrieving all videos by a user
// @Summary Get all videos by user ID
// @Description Get all videos uploaded by a specific user
// @Tags videos
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} models.Video
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{userID}/videos [get]
func (h *VideoHandler) GetUserVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	videos, err := h.videoService.GetUserVideos(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// UpdateVideo handles updating an existing video
// @Summary Update a video
// @Description Update an existing video's information
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param video body models.Video true "Updated video object"
// @Success 200 {object} models.Video
// @Failure 400 {string} string "Invalid video ID or data"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id} [put]
func (h *VideoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	var video models.Video
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	video.ID = id
	if err := h.videoService.UpdateVideo(&video); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

// DeleteVideo handles removing a video
// @Summary Delete a video
// @Description Delete an existing video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid video ID"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id} [delete]
func (h *VideoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	if err := h.videoService.DeleteVideo(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListVideos handles retrieving a list of videos with pagination
// @Summary List all videos
// @Description Get a paginated list of all videos
// @Tags videos
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {array} models.Video
// @Failure 500 {string} string "Internal server error"
// @Router /videos [get]
func (h *VideoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	videos, err := h.videoService.ListVideos(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// RegisterRoutes registers the video routes
func (h *VideoHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/videos", h.CreateVideo).Methods("POST")
	r.HandleFunc("/videos/{id}", h.GetVideo).Methods("GET")
	r.HandleFunc("/videos/{id}", h.UpdateVideo).Methods("PUT")
	r.HandleFunc("/videos/{id}", h.DeleteVideo).Methods("DELETE")
	r.HandleFunc("/videos", h.ListVideos).Methods("GET")
	r.HandleFunc("/users/{userID}/videos", h.GetUserVideos).Methods("GET")
}
