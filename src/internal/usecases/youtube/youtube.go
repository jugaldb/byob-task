package youtube_usecase

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	errorsDom "jugaldb.com/byob_task/src/internal/domain/errors"
	youtube_domain "jugaldb.com/byob_task/src/internal/domain/youtube"
	"jugaldb.com/byob_task/src/utils"
	"regexp"
	"strings"
	"sync"
)

type youtubeUseCase struct {
	config *utils.Config
	logger *utils.StandardLogger
}

type UseCase interface {
	GetOne(ctx context.Context, query string) (*youtube_domain.ChannelMetrics, error)
	GetBatch(ctx context.Context, queries *youtube_domain.BatchRequest) ([]*youtube_domain.ChannelData, error)
}

func (y *youtubeUseCase) GetOne(ctx context.Context, query string) (*youtube_domain.ChannelMetrics, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(y.config.YoutubeAPIKey))
	if err != nil {
		return nil, err
	}

	var channelID string
	match, err := regexp.MatchString(`^(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\w+\/|(?:user|channel)\/|[^\/\n\s]+\/?$))([^\/\n\s]+)`, query)
	if err != nil {
		return nil, err
	}
	if match {
		// Extract channel ID from channel URL
		channelID = regexp.MustCompile(`[^\/\n\s]+`).FindString(query)
	} else {
		// Search for the channel by name
		part := []string{"id"}
		searchResponse, err := service.Search.List(part).Q(query).Type("channel").MaxResults(1).Do()
		if err != nil {
			return nil, errorsDom.APIError("Error " + err.Error())
		}
		if len(searchResponse.Items) == 0 {
			return nil, errorsDom.APIError("no channel found :  " + query)
		}
		channelID = searchResponse.Items[0].Id.ChannelId
	}

	// Fetch channel details
	part := []string{"statistics", "snippet", "brandingSettings"}
	channelResponse, err := service.Channels.List(part).Id(channelID).Do()
	if err != nil {
		return nil, errorsDom.APIError("Error  " + err.Error())
	}

	if len(channelResponse.Items) == 0 {
		return nil, errorsDom.APIError("Channel not found: " + channelID)
	}

	var channelDescription = channelResponse.Items[0].Snippet.Description

	// Extract keywords from channel description
	keywords := extractKeywords(channelDescription)
	country := channelResponse.Items[0].Snippet.Country

	// Extract channel statistics
	channelStats := channelResponse.Items[0].Statistics
	subscriberCount := channelStats.SubscriberCount
	totalViews := channelStats.ViewCount
	var totalVideos = channelStats.VideoCount
	var avgViewsPerVideo float64
	if totalVideos != 0 {
		avgViewsPerVideo = float64(totalViews) / float64(totalVideos)
	}
	profileImageURL := channelResponse.Items[0].Snippet.Thumbnails.High.Url
	bannerImageURL := channelResponse.Items[0].BrandingSettings.Image.BannerExternalUrl
	part = []string{"snippet"}
	playlistItemsResponse, err := service.PlaylistItems.List(part).PlaylistId("UU" + channelID[2:]).MaxResults(1).Do()
	if err != nil {
		return nil, err
	}

	if len(playlistItemsResponse.Items) == 0 {
		return nil, err
	}

	latestVideoID := playlistItemsResponse.Items[0].Snippet.ResourceId.VideoId
	latestVideoURL := "https://www.youtube.com/embed/" + latestVideoID

	metrics := &youtube_domain.ChannelMetrics{
		SubscriberCount:  subscriberCount,
		TotalViews:       totalViews,
		TotalVideos:      totalVideos,
		Country:          country,
		AvgViewsPerVideo: avgViewsPerVideo,
		Keywords:         keywords,
		ChannelName:      channelResponse.Items[0].Snippet.Title,
		ChannelCreated:   channelResponse.Items[0].Snippet.PublishedAt,
		ProfileImageURL:  profileImageURL,
		BannerImageURL:   bannerImageURL,
		LatestVideo:      latestVideoURL,
	}

	return metrics, nil
}

// FetchVideoMetricsBatch fetches metrics for multiple YouTube channels using batch processing
func (y *youtubeUseCase) GetBatch(ctx context.Context, body *youtube_domain.BatchRequest) ([]*youtube_domain.ChannelData, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	queries := body.Channels
	channelDataList := make([]*youtube_domain.ChannelData, len(queries))

	for i, query := range queries {
		wg.Add(1)
		go func(index int, q string) {
			defer wg.Done()
			videoMetrics, err := y.GetOne(ctx, q)
			mu.Lock()
			defer mu.Unlock()
			channelDataList[index] = &youtube_domain.ChannelData{
				Query:        q,
				VideoMetrics: videoMetrics,
				Error:        err,
			}
		}(i, query)
	}

	wg.Wait()
	return channelDataList, nil
}

func extractKeywords(description string) []string {
	// Convert the description to lowercase
	description = strings.ToLower(description)

	var stopWords = map[string]bool{
		"i": true, "me": true, "my": true, "myself": true, "we": true, "our": true, "ours": true, "ourselves": true,
		"you": true, "your": true, "yours": true, "yourself": true, "yourselves": true, "he": true, "him": true,
		"his": true, "himself": true, "she": true, "her": true, "hers": true, "herself": true, "it": true, "its": true,
		"itself": true, "they": true, "them": true, "their": true, "theirs": true, "themselves": true, "what": true,
		"which": true, "who": true, "whom": true, "this": true, "that": true, "these": true, "those": true, "am": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true, "being": true, "have": true,
		"has": true, "had": true, "having": true, "do": true, "does": true, "did": true, "doing": true, "a": true,
		"an": true, "the": true, "and": true, "but": true, "if": true, "or": true, "because": true, "as": true, "until": true,
		"while": true, "of": true, "at": true, "by": true, "for": true, "with": true, "about": true, "against": true,
		"between": true, "into": true, "through": true, "during": true, "before": true, "after": true, "above": true,
		"below": true, "to": true, "from": true, "up": true, "down": true, "in": true, "out": true, "on": true, "off": true,
		"over": true, "under": true, "again": true, "further": true, "then": true, "once": true, "here": true, "there": true,
		"when": true, "where": true, "why": true, "how": true, "all": true, "any": true, "both": true, "each": true, "few": true,
		"more": true, "most": true, "other": true, "some": true, "such": true, "no": true, "nor": true, "not": true, "only": true,
		"own": true, "same": true, "so": true, "than": true, "too": true, "very": true, "s": true, "t": true, "can": true, "will": true,
		"just": true, "don": true, "should": true, "now": true, "d": true, "ll": true, "m": true, "o": true, "re": true, "ve": true,
		"y": true, "ain": true, "aren": true, "couldn": true, "didn": true, "doesn": true, "hadn": true, "hasn": true, "haven": true,
		"isn": true, "ma": true, "mightn": true, "mustn": true, "needn": true, "shan": true, "shouldn": true, "wasn": true, "weren": true,
		"won": true, "wouldn": true,
	}

	// Tokenize the description
	words := strings.Fields(description)

	// Remove stop words and non-alphabetic characters
	var keywords []string
	for _, word := range words {
		if !stopWords[word] && isAlphabetic(word) {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func New(config *utils.Config, logger *utils.StandardLogger) UseCase {
	youtubeUsecaseObj := &youtubeUseCase{
		config: config,
		logger: logger,
	}

	return youtubeUsecaseObj
}
