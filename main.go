package main

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"sealway-strava/api"
	"sealway-strava/graph"
	"sealway-strava/infra"
	"sealway-strava/strava"
	"time"
)

var applicationSlug = "integration-strava"

// ENVs
var connectionString = infra.EnvOrDefault("SEALWAY_ConnectionStrings__Mongo__Connection", "mongodb://localhost:27017")
var stravaClientId = os.Getenv("SEALWAY_Services__Strava__Client")
var stravaSecretId = os.Getenv("SEALWAY_Services__Strava__Secret")
var port = infra.EnvOrDefault("PORT", "8080")

func main() {
	stravaClient := strava.NewAPIClient(strava.NewConfiguration())

	ctx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	err, stravaRepository := api.InitStravaRepository(connectionString, ctx)
	if err != nil {
		panic(err)
	}

	var stravaService = &api.StravaService{
		ClientId:     stravaClientId,
		SecretId:     stravaSecretId,
		StravaClient: stravaClient,
	}

	var backgroundWorker = &BackgroundWorker{
		StravaService:    stravaService,
		StravaRepository: stravaRepository,
	}
	queue := backgroundWorker.RunBackgroundWorker()

	router := mux.NewRouter()
	defaultApi := &api.DefaultApi{
		Router:          router,
		ApplicationSlug: applicationSlug,
	}

	var restApi = &api.SubscriptionApi{
		Queue:      queue,
		DefaultApi: defaultApi,
	}
	restApi.RegisterApiRoutes()

	graphqlApi := graph.GraphqlApi{
		Resolvers: &graph.Resolver{
			Repository: stravaRepository,
		},
		DefaultApi: defaultApi,
	}
	graphqlApi.RegisterGraphQl()

	apiServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	var graceful = &infra.Graceful{
		StartAction: func() error {
			return apiServer.ListenAndServe()
		},
		DeferAction: func(ctx context.Context) error {
			close(queue)
			cancelMongo()
			stravaRepository.Client.Disconnect(ctx)

			return nil
		},
		ShutdownAction: func(ctx context.Context) error {
			return apiServer.Shutdown(ctx)
		},
	}

	graceful.RunAndWait()
}
