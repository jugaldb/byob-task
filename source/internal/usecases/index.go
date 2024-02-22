package usecases

import (
	youtube_usecase "jugaldb.com/byob_task/src/internal/usecases/youtube"
	"jugaldb.com/byob_task/src/utils"
)

type UseCases struct {
	Youtube youtube_usecase.UseCase
}

func InitUseCases(config *utils.Config, logger *utils.StandardLogger) *UseCases {
	youtubeUsecase := youtube_usecase.New(config, logger)
	return &UseCases{
		Youtube: youtubeUsecase,
	}
}
