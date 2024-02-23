package youtube_usecase

import (
	"context"
	"fmt"
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
			return nil, errorsDom.APIError("no channel found with the name " + query)
		}
		if len(searchResponse.Items) == 0 {
			return nil, errorsDom.APIError("no channel found with the name " + query)
		}
		channelID = searchResponse.Items[0].Id.ChannelId
	}

	// Fetch channel details
	part := []string{"statistics", "snippet"}
	channelResponse, err := service.Channels.List(part).Id(channelID).Do()
	if err != nil {
		return nil, errorsDom.APIError("no response found with the name " + query)
	}

	if len(channelResponse.Items) == 0 {
		return nil, errorsDom.APIError("channel with ID '%s' not found" + channelID)
	}

	var channelDescription = channelResponse.Items[0].Snippet.Description
	fmt.Println(channelDescription)

	// Extract keywords from channel description
	keywords := extractKeywords(channelDescription)

	// Extract channel statistics
	channelStats := channelResponse.Items[0].Statistics
	subscriberCount := channelStats.SubscriberCount
	totalViews := channelStats.ViewCount

	part = []string{"contentDetails"}
	// Fetch playlist details to get the number of videos
	playlistResponse, err := service.Playlists.List(part).ChannelId(channelID).Do()
	if err != nil {
		return nil, errorsDom.APIError("channel with ID '%s' not found" + channelID)
	}

	var totalVideos uint64
	for _, playlist := range playlistResponse.Items {
		totalVideos += uint64(playlist.ContentDetails.ItemCount)
	}

	// Fetch videos from uploads playlist to calculate engagement metrics
	uploadsPlaylistID := playlistResponse.Items[0].Id
	playlistItemsResponse, err := service.PlaylistItems.List(part).PlaylistId(uploadsPlaylistID).Do()
	if err != nil {
		return nil, err
	}

	var totalLikes, totalComments uint64
	for _, playlistItem := range playlistItemsResponse.Items {
		videoID := playlistItem.ContentDetails.VideoId
		part = []string{"statistics"}
		videoResponse, err := service.Videos.List(part).Id(videoID).Do()
		if err != nil {
			return nil, errorsDom.APIError("channel with ID '%s' not found" + channelID)
		}
		videoStats := videoResponse.Items[0].Statistics
		likes := videoStats.LikeCount
		comments := videoStats.CommentCount
		totalLikes += uint64(likes)
		totalComments += uint64(comments)
	}

	var avgViewsPerVideo, avgLikesPerVideo, avgCommentsPerVideo float64
	if totalVideos != 0 {
		avgViewsPerVideo = float64(totalViews) / float64(totalVideos)
		avgLikesPerVideo = float64(totalLikes) / float64(totalVideos)
		avgCommentsPerVideo = float64(totalComments) / float64(totalVideos)
	}

	metrics := &youtube_domain.ChannelMetrics{
		SubscriberCount:     subscriberCount,
		TotalViews:          totalViews,
		TotalVideos:         totalVideos,
		TotalLikes:          totalLikes,
		TotalComments:       totalComments,
		AvgViewsPerVideo:    avgViewsPerVideo,
		AvgLikesPerVideo:    avgLikesPerVideo,
		AvgCommentsPerVideo: avgCommentsPerVideo,
		Keywords:            keywords,
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
