package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"jugaldb.com/byob_task/src/delivery/http/common"
	errorsDom "jugaldb.com/byob_task/src/internal/domain/errors"
	youtube_domain "jugaldb.com/byob_task/src/internal/domain/youtube"
	youtube_usecase "jugaldb.com/byob_task/src/internal/usecases/youtube"
	"jugaldb.com/byob_task/src/utils"
	"net/http"
	"time"
)

type YoutubeController interface {
	GetOne(w http.ResponseWriter, r *http.Request)
	GetBatch(w http.ResponseWriter, r *http.Request)
	Home(w http.ResponseWriter, r *http.Request)
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
		err := errorsDom.InvalidBody("channel missing")
		common.HandleError(r.Context(), w, err)
		return
	} else {
		channelLink = channelLinkQuery[0]
	}
	channelDetails, err := y.youtubeUsecase.GetOne(r.Context(), channelLink)
	fmt.Println("channelDetails", channelDetails)
	fmt.Println(err)
	if err != nil {
		common.HandleError(r.Context(), w, err)
		return
	} else {
		common.SendJson(w, channelDetails)
	}
}

func (y *youtubeController) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	common.SendJson(w, map[string]any{
		"hello":        "world",
		"current_time": time.Now(),
		"link":         "https://github.com/jugaldb/byob-task",
	})
}

func (y *youtubeController) GetBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = errorsDom.InternalServerError(err)
		common.HandleError(r.Context(), w, err)
		return
	}
	body := &youtube_domain.BatchRequest{}
	err = json.Unmarshal(b, body)
	if err != nil {
		err = errorsDom.InternalServerError(err)
		common.HandleError(r.Context(), w, err)
		return
	}
	channelDetails, err := y.youtubeUsecase.GetBatch(r.Context(), body)
	if err != nil {
		common.HandleError(r.Context(), w, err)
		return
	} else {
		common.SendJson(w, channelDetails)
	}
}

func NewYoutubeController(config *utils.Config, youtubeUsecase youtube_usecase.UseCase) YoutubeController {
	return &youtubeController{
		config:         config,
		youtubeUsecase: youtubeUsecase,
	}
}
