package internal

import (
	"context"
	"net/http"
	"sealway-strava/domain"
	"sealway-strava/interfaces/graph"
	"sealway-strava/interfaces/rest"
	"sealway-strava/pkg/graceful"
	usercase "sealway-strava/usecase"
)

type ServerConfig struct {
	Port string
}

type SealwayStravaApp struct {
	*graceful.Graceful
}

func MakeApplication(config *ServerConfig, apiConfig *rest.ApiConfig, queue *domain.ActivitiesQueue, manager *usercase.SubscriptionManager, service *usercase.StravaService, restApi *rest.SubscriptionApi, graphApi *graph.GraphqlApi) graceful.Application {
	restApi.RegisterHealth()
	restApi.RegisterApiRoutes()
	graphApi.RegisterGraphQl()

	apiServer := &http.Server{
		Addr:    ":" + config.Port,
		Handler: apiConfig.Router,
	}

	var graceful = &graceful.Graceful{
		StartAction: func() error {
			return apiServer.ListenAndServe()
		},
		DeferAction: func(ctx context.Context) error {
			queue.Close()
			manager.Close()

			return service.Close()
		},
		ShutdownAction: func(ctx context.Context) error {
			return apiServer.Shutdown(ctx)
		},
	}

	return &SealwayStravaApp{
		Graceful: graceful,
	}
}
