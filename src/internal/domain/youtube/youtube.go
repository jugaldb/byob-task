package youtube_domain

type ChannelMetrics struct {
	SubscriberCount  uint64
	TotalViews       uint64
	TotalVideos      uint64
	Country          string
	AvgViewsPerVideo float64
	Keywords         []string
	ChannelCreated   string
	ChannelName      string
	ProfileImageURL  string
	BannerImageURL   string
	LatestVideo      string
}

type ChannelData struct {
	Query        string
	VideoMetrics *ChannelMetrics
	Error        error
}

type BatchRequest struct {
	Channels []string `json:"channels"`
}
