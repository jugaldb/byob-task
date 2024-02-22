package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	httpDelivery "jugaldb.com/byob_task/src/delivery/http"
	httpMiddleware "jugaldb.com/byob_task/src/delivery/middleware"
	"jugaldb.com/byob_task/src/internal/usecases"
	"jugaldb.com/byob_task/src/utils"
	"jugaldb.com/byob_task/src/utils/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := utils.NewConfig()
	logger := utils.NewLogger(config)
	utils.SetAppLogger(logger)
	//globalContext := context.Background()
	metrics.NewMetricCounters()

	// Sentry
	if config.DGN != "local" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.SentryDSN,
			Debug:            config.Debug,
			Environment:      config.DGN,
			IgnoreErrors:     []string{},
			EnableTracing:    true,
			TracesSampleRate: 0.1,
		})
		logger.Debugf("Sentry initiated" + config.SentryDSN)
		defer sentry.Flush(config.SentryContextTimeout)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if config.DGN == "STAG" || config.DGN == "PROD" {
		//redisOpt.Password = config.RedisPassword
		//redisOpt.TLSConfig = &tls.Config{
		//	MinVersion: tls.VersionTLS12,
		//	ServerName: config.RedisTLSDomain,
		//}
	}

	useCases := usecases.InitUseCases(config, logger)
	middlewares := httpMiddleware.New(config, useCases.Youtube)

	logger.Debug("Intitializing Metrics Server...")
	go metrics.RunMetricsServer(config)
	_, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), httpMiddleware.ApplyMiddlewares(http.DefaultServeMux, middlewares, config)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	httpDelivery.NewRestDelivery(config, useCases)

	<-done
	log.Print("OS kill received")

	defer func() {
		cancel()
	}()
	log.Print("Server Exited Properly")
}
