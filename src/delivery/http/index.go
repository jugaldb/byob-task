package httpDelivery

import (
	"github.com/gorilla/mux"
	"jugaldb.com/byob_task/src/delivery/controllers/youtube"
	"jugaldb.com/byob_task/src/delivery/routers"
	"jugaldb.com/byob_task/src/internal/usecases"
	"jugaldb.com/byob_task/src/utils"
	"net/http"
)

func NewRestDelivery(config *utils.Config, useCases *usecases.UseCases) {
	r := mux.NewRouter()
	routers.SetMetricsRoute(r)
	routers.SetYoutubeRoutes(youtube.NewYoutubeController(config, useCases.Youtube), r)
	http.Handle("/", r)
	utils.GetAppLogger().Infof("Rest delivery listening on port %d", config.ServerPort)
}
