package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalAPIRequests *prometheus.CounterVec
var TotalEventsReceived *prometheus.CounterVec
var TotalEventsProcessed *prometheus.CounterVec
var TotalEventsWritten *prometheus.CounterVec
var TotalEventsDropped *prometheus.CounterVec
var EventProcessingTime *prometheus.HistogramVec

var apiMetricPrefic = "halfblood_api"
var eventMetricPrefix = "events"

func NewMetricCounters() {
	TotalAPIRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName(apiMetricPrefic, "received_total"),
		Help: "total number of events received",
	}, []string{"api_type", "endpoint"})

	TotalEventsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName(eventMetricPrefix, "received_total"),
		Help: "total number of events received",
	}, []string{"event_type", "entity"})

	TotalEventsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName(eventMetricPrefix, "processed_total"),
		Help: "total number of events processed",
	}, []string{"event_type", "entity"})

	TotalEventsWritten = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName(eventMetricPrefix, "written_total"),
		Help: "total number of events written",
	}, []string{"event_type", "entity"})

	TotalEventsDropped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName(eventMetricPrefix, "dropped_total"),
		Help: "total number of events dropped",
	}, []string{"event_type", "entity", "error"})

	EventProcessingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: metricName(eventMetricPrefix, "processing_time_seconds"),
		Help: "time taken to process an event",
	}, []string{"event_type", "entity"})
}

func metricName(prefix, name string) string {
	return fmt.Sprintf("%s_%s", prefix, name)
}
