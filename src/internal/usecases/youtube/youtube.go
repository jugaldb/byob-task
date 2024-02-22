package youtube_usecase

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	youtube_domain "jugaldb.com/byob_task/src/internal/domain/youtube"
	"jugaldb.com/byob_task/src/utils"
	"regexp"
)

type youtubeUseCase struct {
	config *utils.Config
	logger *utils.StandardLogger
}

type UseCase interface {
	GetOne(ctx context.Context, query string) (*youtube_domain.ChannelMetrics, error)
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
			return nil, err
		}
		if len(searchResponse.Items) == 0 {
			return nil, fmt.Errorf("no channel found with the name '%s'", query)
		}
		channelID = searchResponse.Items[0].Id.ChannelId
	}

	// Fetch channel details
	part := []string{"statistics"}
	channelResponse, err := service.Channels.List(part).Id(channelID).Do()
	if err != nil {
		return nil, err
	}

	if len(channelResponse.Items) == 0 {
		return nil, fmt.Errorf("channel with ID '%s' not found", channelID)
	}

	// Extract channel statistics
	channelStats := channelResponse.Items[0].Statistics
	subscriberCount := channelStats.SubscriberCount
	totalViews := channelStats.ViewCount

	part = []string{"contentDetails"}
	// Fetch playlist details to get the number of videos
	playlistResponse, err := service.Playlists.List(part).ChannelId(channelID).Do()
	if err != nil {
		return nil, err
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
			return nil, err
		}
		videoStats := videoResponse.Items[0].Statistics
		likes := videoStats.LikeCount
		comments := videoStats.CommentCount
		totalLikes += uint64(likes)
		totalComments += uint64(comments)
	}

	metrics := &youtube_domain.ChannelMetrics{
		SubscriberCount: subscriberCount,
		TotalViews:      totalViews,
		TotalVideos:     totalVideos,
		TotalLikes:      totalLikes,
		TotalComments:   totalComments,
	}

	return metrics, nil
}

func New(config *utils.Config, logger *utils.StandardLogger) UseCase {
	youtubeUsecaseObj := &youtubeUseCase{
		config: config,
		logger: logger,
	}

	return youtubeUsecaseObj
}
