package routers

import (
	"github.com/gorilla/mux"
	"jugaldb.com/byob_task/src/delivery/controllers/youtube"
)

type YoutubeController = youtube.YoutubeController

func SetYoutubeRoutes(youtubeController YoutubeController, r *mux.Router) {
	youtubeRoutes := r.PathPrefix("/youtube").Subrouter()
	youtubeRoutes.HandleFunc("/getOne", youtubeController.GetOne).Methods("GET")
	youtubeRoutes.HandleFunc("/getBatch", youtubeController.GetBatch).Methods("POST")
}
