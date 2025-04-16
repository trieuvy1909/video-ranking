package ws

// UpdateType represents the type of update
type UpdateType string

const (
	VideoScoreUpdate     UpdateType = "video_score"
	VideoStatsUpdate     UpdateType = "video_stats"
	InteractionUpdate    UpdateType = "interaction"
	TrendingVideosUpdate UpdateType = "trending_videos"
)

// Update represents a WebSocket update message
type Update struct {
	Type    UpdateType  `json:"type"`
	Payload interface{} `json:"payload"`
}

// VideoScorePayload represents the payload for video score updates
type VideoScorePayload struct {
	VideoID uint    `json:"video_id"`
	Score   float64 `json:"score"`
}

// VideoStatsPayload represents the payload for video stats updates
type VideoStatsPayload struct {
	VideoID  uint  `json:"video_id"`
	Views    int64 `json:"views,omitempty"`
	Likes    int64 `json:"likes,omitempty"`
	Dislikes int64 `json:"dislikes,omitempty"`
	Comments int64 `json:"comments,omitempty"`
}

// InteractionPayload represents the payload for interaction updates
type InteractionPayload struct {
	VideoID uint   `json:"video_id"`
	UserID  uint   `json:"user_id"`
	Type    string `json:"type"`
}

// TrendingVideosPayload represents the payload for trending videos updates
type TrendingVideosPayload struct {
	Videos []interface{} `json:"videos"`
}
