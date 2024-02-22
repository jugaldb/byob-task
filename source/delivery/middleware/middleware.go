package httpMiddleware

import (
	youtube_usecase "jugaldb.com/byob_task/src/internal/usecases/youtube"
	"jugaldb.com/byob_task/src/utils"
)

type Middlewares struct {
	Config         *utils.Config
	youtubeUseCase youtube_usecase.UseCase
}

func New(config *utils.Config, youtubeUseCase youtube_usecase.UseCase) *Middlewares {
	return &Middlewares{
		config,
		youtubeUseCase,
	}
}
