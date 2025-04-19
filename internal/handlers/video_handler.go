package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"strings"
	"github.com/go-playground/validator/v10"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/trieuvy/video-ranking/internal/params/request"
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
// @Param video body request.Video true "Video object"
// @Success 200 {object} request.Video
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /videos [post]
func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var video request.Video
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	err := validate.Struct(video)
	if err != nil {
		var sb strings.Builder
		for _, e := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on the '%s' rule\n", e.Field(), e.Tag()))
		}
		http.Error(w, sb.String(), http.StatusBadRequest)
		return
	}
	videoModel := models.Video{
		Title: video.Title,
		Description: video.Description,
		Views: 0,
		Likes: 0,
		Comments: 0,
	}
	if err := h.videoService.CreateVideo(&videoModel); err != nil {
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
// @Success 200 {object} request.Video
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

// UpdateVideo handles updating an existing video
// @Summary Update a video
// @Description Update an existing video's information
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param video body request.Video true "Updated video object"
// @Success 200 {object} request.Video
// @Failure 400 {string} string "Invalid video ID or data"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id} [put]
func (h *VideoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	var video request.Video
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	err = validate.Struct(video)
	if err != nil {
		var sb strings.Builder
		for _, e := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on the '%s' rule\n", e.Field(), e.Tag()))
		}
		http.Error(w, sb.String(), http.StatusBadRequest)
		return
	}
	videoModel, err := h.videoService.GetVideo(id)
	if err!= nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	videoModel.Title = video.Title
	videoModel.Description = video.Description
	if err := h.videoService.UpdateVideo(videoModel); err != nil {
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
// @Success 200 {array} request.Video
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
}

// ChangeLikesAmount handles changing the amount of likes for a video
// @Summary Change likes amount
// @Description Change the number of likes for a specific video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param step query int true "Step"
// @Success 200 {string} string "Likes amount changed"
// @Failure 400 {string} string "Invalid video ID or step"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id}/likes [patch]
func (h *VideoHandler) ChangeLikesAmount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	step, err := strconv.Atoi(r.URL.Query().Get("step"))
	if err != nil {
		http.Error(w, "Invalid step", http.StatusBadRequest)
		return
	}

	if err := h.videoService.ChangeLikesAmount(id, step); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Likes amount changed"))
}

// ChangeViewsAmount handles changing the amount of views for a video
// @Summary Change views amount
// @Description Change the number of views for a specific video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param step query int true "Step"
// @Success 200 {string} string "Views amount changed"
// @Failure 400 {string} string "Invalid video ID or step"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id}/views [patch]
func (h *VideoHandler) ChangeViewsAmount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	step, err := strconv.Atoi(r.URL.Query().Get("step"))
	if err != nil {
		http.Error(w, "Invalid step", http.StatusBadRequest)
		return
	}

	if err := h.videoService.ChangeViewsAmount(id, step); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Views amount changed"))
}

// ChangeCommentsAmount handles changing the amount of comments for a video
// @Summary Change comments amount
// @Description Change the number of comments for a specific video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param step query int true "Step"
// @Success 200 {string} string "Comments amount changed"
// @Failure 400 {string} string "Invalid video ID or step"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/{id}/comments [patch]
func (h *VideoHandler) ChangeCommentsAmount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	step, err := strconv.Atoi(r.URL.Query().Get("step"))
	if err != nil {
		http.Error(w, "Invalid step", http.StatusBadRequest)
		return
	}

	if err := h.videoService.ChangeCommentsAmount(id, step); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comments amount changed"))
}
