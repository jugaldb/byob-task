package metrics

import (
	"jugaldb.com/byob_task/src/utils"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunMetricsServer(config *utils.Config) {
	http.Handle("/metrics", promhttp.Handler())
	utils.GetAppLogger().Debugf(":" + config.MetricsPort)
	http.ListenAndServe(":"+config.MetricsPort, nil)
}
