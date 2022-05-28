package internal

import (
	"context"
	"net/http"
	"sealway-strava/interfaces/graph"
	"sealway-strava/interfaces/rest"
	"sealway-strava/pkg/closer"
	"sealway-strava/pkg/graceful"
)

type ServerConfig struct {
	Port string
}

type SealwayStravaApp struct {
	*graceful.Graceful
}

func MakeApplication(collection *closer.CloserCollection, config *ServerConfig, apiConfig *rest.ApiConfig, restApi *rest.SubscriptionApi, graphApi *graph.GraphqlApi) graceful.Application {
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
			return collection.Close(ctx)
		},
		ShutdownAction: func(ctx context.Context) error {
			return apiServer.Shutdown(ctx)
		},
	}

	return &SealwayStravaApp{
		Graceful: graceful,
	}
}
