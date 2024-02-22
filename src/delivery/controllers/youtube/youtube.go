package youtube

import (
	"jugaldb.com/byob_task/src/delivery/http/common"
	youtube_usecase "jugaldb.com/byob_task/src/internal/usecases/youtube"
	"jugaldb.com/byob_task/src/utils"
	"net/http"
)

type YoutubeController interface {
	GetOne(w http.ResponseWriter, r *http.Request)
}

type youtubeController struct {
	config         *utils.Config
	youtubeUsecase youtube_usecase.UseCase
}

func (y *youtubeController) GetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	channelLinkQuery := r.URL.Query()["channelLink"]
	var channelLink string
	if len(channelLinkQuery) == 0 {
		channelLink = "0"
	} else {
		channelLink = channelLinkQuery[0]
	}
	channelDetails, err := y.youtubeUsecase.GetOne(r.Context(), channelLink)
	if err != nil {
		common.HandleError(r.Context(), w, err)
		return
	} else {
		common.SendJson(w, map[string]any{
			"data": channelDetails,
		})
	}
}

func NewYoutubeController(config *utils.Config, youtubeUsecase youtube_usecase.UseCase) YoutubeController {
	return &youtubeController{
		config:         config,
		youtubeUsecase: youtubeUsecase,
	}
}
