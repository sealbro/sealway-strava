package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/dig"
	"os"
	"sealway-strava/infrastructure"
	"sealway-strava/interfaces/graph"
	"sealway-strava/interfaces/rest"
	"sealway-strava/internal"
	"sealway-strava/pkg/closer"
	"sealway-strava/pkg/env"
	"sealway-strava/pkg/graceful"
	"sealway-strava/pkg/logger"
	"sealway-strava/repository"
	"sealway-strava/usecase"
	"strconv"
	"time"
)

var connectionString = env.EnvOrDefault("MONGO_CONNECTION", "mongodb://localhost:27017")
var stravaClientId = os.Getenv("STRAVA_CLIENT")
var stravaSecretId = os.Getenv("STRAVA_SECRET")

var activityBatchSize, _ = strconv.Atoi(env.EnvOrDefault("ACTIVITY_BATCH_SIZE", "50"))
var activityBatchTime, _ = time.ParseDuration(env.EnvOrDefault("ACTIVITY_BATCH_TIME", "45s"))

var port = env.EnvOrDefault("PORT", "8080")
var applicationSlug = env.EnvOrDefault("SLUG", "integration-strava")

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

	panicOrProvide(container, func() *usercase.BatchConfig {
		return &usercase.BatchConfig{
			ActivityBatchSize: activityBatchSize,
			ActivityBatchTime: activityBatchTime,
		}
	})
	panicOrProvide(container, func() *rest.ApiConfig {
		return &rest.ApiConfig{
			ApplicationSlug: applicationSlug,
			Router:          mux.NewRouter(),
		}
	})
	panicOrProvide(container, func() *internal.ServerConfig {
		return &internal.ServerConfig{
			Port: port,
		}
	})
	panicOrProvide(container, func() *infrastructure.StravaConfig {
		return &infrastructure.StravaConfig{
			ClientId: stravaClientId,
			SecretId: stravaSecretId,
		}
	})
	panicOrProvide(container, func() *repository.MongoConfig {
		return &repository.MongoConfig{
			ConnectionString: connectionString,
		}
	})

	panicOrProvide(container, usercase.MakeSubscriptionManager)
	panicOrProvide(container, usercase.MakeStravaService)
	panicOrProvide(container, infrastructure.MakeStravaClient)
	panicOrProvide(container, repository.MakeStravaRepository)

	panicOrProvide(container, usercase.MakeBackgroundWorker)
	panicOrProvide(container, rest.MakeRestApi)
	panicOrProvide(container, rest.MakeSubscriptionApi)
	panicOrProvide(container, graph.MakeGraphqlApi)
	panicOrProvide(container, closer.MakeCloserCollection)
	panicOrProvide(container, internal.MakeApplication)

	return container
}

func panicOrProvide(container *dig.Container, constructor interface{}, opts ...dig.ProvideOption) {
	err := container.Provide(constructor, opts...)
	if err != nil {
		logger.Fatalf("container.Provide: $v", err)
	}
}
