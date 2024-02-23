package youtube_domain

type ChannelMetrics struct {
	SubscriberCount     uint64
	TotalViews          uint64
	TotalVideos         uint64
	TotalLikes          uint64
	TotalComments       uint64
	AvgViewsPerVideo    float64
	AvgLikesPerVideo    float64
	AvgCommentsPerVideo float64
	Keywords            []string
}

type ChannelData struct {
	Query        string
	VideoMetrics *ChannelMetrics
	Error        error
}

type BatchRequest struct {
	Channels []string `json:"channels"`
}
