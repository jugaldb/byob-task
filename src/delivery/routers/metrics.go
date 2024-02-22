package routers

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetMetricsRoute(r *mux.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
