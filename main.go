package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/dig"
	"os"
	"sealway-strava/infrastructure"
	"sealway-strava/interfaces/graph"
	"sealway-strava/interfaces/rest"
	"sealway-strava/internal"
	"sealway-strava/pkg/graceful"
	"sealway-strava/repository"
	"sealway-strava/usecase"
	"strconv"
	"time"
)

var connectionString = graceful.EnvOrDefault("MONGO_CONNECTION", "mongodb://localhost:27017")
var stravaClientId = os.Getenv("STRAVA_CLIENT")
var stravaSecretId = os.Getenv("STRAVA_SECRET")

var activityBatchSize, _ = strconv.Atoi(graceful.EnvOrDefault("ACTIVITY_BATCH_SIZE", "50"))
var activityBatchTime, _ = time.ParseDuration(graceful.EnvOrDefault("ACTIVITY_BATCH_TIME", "45s"))

var port = graceful.EnvOrDefault("PORT", "8080")
var applicationSlug = graceful.EnvOrDefault("SLUG", "integration-strava")

func main() {
	container := setupContainer()

	err := container.Invoke(func(application graceful.Application) {
		application.RunAndWait()
	})

	if err != nil {
		panic(err)
	}
}

func setupContainer() *dig.Container {
	container := dig.New()

	container.Provide(func() *usercase.BatchConfig {
		return &usercase.BatchConfig{
			ActivityBatchSize: activityBatchSize,
			ActivityBatchTime: activityBatchTime,
		}
	})
	container.Provide(func() *rest.ApiConfig {
		return &rest.ApiConfig{
			ApplicationSlug: applicationSlug,
			Router:          mux.NewRouter(),
		}
	})
	container.Provide(func() *internal.ServerConfig {
		return &internal.ServerConfig{
			Port: port,
		}
	})
	container.Provide(func() *infrastructure.StravaConfig {
		return &infrastructure.StravaConfig{
			ClientId: stravaClientId,
			SecretId: stravaSecretId,
		}
	})
	container.Provide(func() *repository.MongoConfig {
		return &repository.MongoConfig{
			ConnectionString: connectionString,
		}
	})

	container.Provide(usercase.MakeSubscriptionManager)
	container.Provide(usercase.MakeStravaService)
	container.Provide(infrastructure.MakeStravaClient)
	err := container.Provide(repository.MakeStravaRepository)
	if err != nil {
		panic(err)
	}

	container.Provide(usercase.MakeBackgroundWorker)
	container.Provide(rest.MakeRestApi)
	container.Provide(rest.MakeSubscriptionApi)
	container.Provide(graph.MakeGraphqlApi)
	container.Provide(internal.MakeApplication)

	return container
}
