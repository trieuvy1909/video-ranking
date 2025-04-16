package services

import (
    "context"
    "log"
    "time"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
    "github.com/trieuvy/video-ranking/internal/models"
    "github.com/trieuvy/video-ranking/internal/repositories"
)

// VideoService handles business logic for videos
type VideoService struct {
    repo        *repositories.VideoRepository
    redisClient *redis.Client
}

// NewVideoService creates a new video service
func NewVideoService(repo *repositories.VideoRepository,redisClient *redis.Client ) *VideoService {
    return &VideoService{
        repo: repo,
        redisClient: redisClient,
    }
}

// CreateVideo creates a new video
func (s *VideoService) CreateVideo(video *models.Video) error {
    err := s.repo.Create(video)
    if err != nil {
        return err
    }
    return nil
}

// GetVideo retrieves a video by ID
func (s *VideoService) GetVideo(id uuid.UUID) (*models.Video, error) {
    return s.repo.FindByID(id)
}

// GetUserVideos retrieves all videos by a user
func (s *VideoService) GetUserVideos(userID uuid.UUID) ([]models.Video, error) {
    return s.repo.FindByUser(userID)
}

// UpdateVideo updates an existing video
func (s *VideoService) UpdateVideo(video *models.Video) error {
    return s.repo.Update(video)
}

// DeleteVideo removes a video
func (s *VideoService) DeleteVideo(id uuid.UUID) error {
    err:= s.repo.Delete(id)
    if err != nil {
        return err
    }
    // Remove video from Redis
    ctx := context.Background()
    err = s.redisClient.ZRem(ctx, "video:scores", id.String()).Err()
    if err != nil {
        log.Printf("Error removing video from Redis: %v", err)
        return err
    }
    trendingVideos, err := s.getTop10TrendingVideos(ctx)
    if err != nil {
        log.Printf("Error getting top trending videos: %v", err)
        return err
    }
    
    jsonData,err := PrepareVideoData(trendingVideos)
    if err != nil {
        log.Printf("Error preparing video data: %v", err)
        return err
    }
    return SendNotification(s.redisClient, "trending_videos", string(jsonData))
}

// ListVideos retrieves a list of videos with pagination
func (s *VideoService) ListVideos(page, pageSize int) ([]models.Video, error) {
    offset := (page - 1) * pageSize
    return s.repo.List(offset, pageSize)
}

// UpdateVideoScore updates the score of a video
func (s *VideoService) UpdateVideoScore(id uuid.UUID, score float64) error {
    return s.repo.UpdateScore(id, score)
}

func (s *VideoService) ChangeLikesAmount(videoID uuid.UUID, step int) error {
    err := s.repo.ChangeLikes(videoID,step)
    if err != nil {
        return err
    }
    
    // Update ranking if ranking service is available
    ctx := context.Background()
    return s.UpdateAndNotifyRanking(ctx, videoID)
}

func (s *VideoService) ChangeViewsAmount(videoID uuid.UUID, step int) error {
    err := s.repo.ChangeViews(videoID,step)
    if err != nil {
        return err
    }
    
    // Update ranking if ranking service is available
    ctx := context.Background()
    return s.UpdateAndNotifyRanking(ctx, videoID)
}

func (s *VideoService) ChangeCommentsAmount(videoID uuid.UUID, step int) error {
    err := s.repo.ChangeComments(videoID, step)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    return s.UpdateAndNotifyRanking(ctx, videoID)
}

// UpdateAndNotifyRanking updates the ranking of a video and notifies clients
func (s *VideoService) UpdateAndNotifyRanking(ctx context.Context, videoID uuid.UUID) error {
    
    video, err := s.GetVideo(videoID)
    if err != nil {
        return err
    }
    
    newScore := calculateEngagementScore(video.Views, video.Likes, video.Comments)
    
    err = s.UpdateVideoScore(videoID, newScore)
    if err != nil {
        return err
    }
    
    err = s.redisClient.ZAdd(ctx, "video:scores", &redis.Z{
        Score:  newScore,
        Member: videoID.String(),
    }).Err()
    if err != nil {
        log.Printf("Error updating Redis score: %v", err)
        return err
    }
    
    trendingVideos, err := s.getTop10TrendingVideos(ctx)
    if err != nil {
        log.Printf("Error getting top trending videos: %v", err)
        return err
    }
    
    jsonData,err := PrepareVideoData(trendingVideos)
    if err != nil {
        log.Printf("Error preparing video data: %v", err)
        return err
    }
    return SendNotification(s.redisClient, "trending_videos", string(jsonData))
}

// calculateEngagementScore calculates engagement score
func calculateEngagementScore(views, likes, comments int64) float64 {
    // Weights for different engagement metrics
    const (
        viewWeight    = 1.0
        likeWeight    = 2.0
        commentWeight = 3.0
    )

    // Calculate base engagement score
    engagementScore := float64(views)*viewWeight +
        float64(likes)*likeWeight +
        float64(comments)*commentWeight

    // Normalize score to be positive
    if engagementScore < 0 {
        engagementScore = 0
    }

    return engagementScore
}

func (s *VideoService) getTop10TrendingVideos(ctx context.Context) ([]map[string]interface{}, error) {
    var trendingVideos []map[string]interface{}
    results, err := s.redisClient.ZRevRangeWithScores(ctx, "video:scores", 0, 9).Result()
    if err != nil {
        return trendingVideos,err
    }
    
    if len(results) == 0 {
        return trendingVideos,nil
    }
    
    for i, z := range results {
        videoData := map[string]interface{}{
            "rank":     i + 1,               
            "video_id": z.Member.(string),   
            "score":    z.Score,            
        }
        trendingVideos = append(trendingVideos, videoData)
    }
    return trendingVideos, nil
}

func SendNotification(redisClient *redis.Client, channel string, message string) error {
    ctx := context.Background()
    err := redisClient.Publish(ctx, channel, message).Err()
    if err != nil {
        log.Printf("Error publishing message to Redis: %v", err)
        return err
    }
    return nil
}

func PrepareVideoData(trendingVideos []map[string]interface{}) ([]byte, error) {
    update := map[string]interface{}{
        "type":    "trending_videos",
        "videos":  trendingVideos,
        "updated": time.Now().Format(time.RFC3339),
    }
    return json.Marshal(update)
}
