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
	"strconv"
	"time"
)

// ENVs
var connectionString = infra.EnvOrDefault("SEALWAY_ConnectionStrings__Mongo__Connection", "mongodb://localhost:27017")
var stravaClientId = os.Getenv("SEALWAY_Services__Strava__Client")
var stravaSecretId = os.Getenv("SEALWAY_Services__Strava__Secret")

var activityBatchSize, _ = strconv.Atoi(infra.EnvOrDefault("ACTIVITY_BATCH_SIZE", "50"))
var activityBatchTime, _ = time.ParseDuration(infra.EnvOrDefault("ACTIVITY_BATCH_TIME", "45s"))

var port = infra.EnvOrDefault("PORT", "8080")
var applicationSlug = infra.EnvOrDefault("SLUG", "integration-strava")

func main() {
	stravaClient := strava.NewAPIClient(strava.NewConfiguration())

	ctx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	err, stravaRepository := api.InitStravaRepository(connectionString, ctx)
	if err != nil {
		panic(err)
	}

	subscriptionManager := &graph.SubscriptionManager{
		ActivityBatchSize: activityBatchSize,
		ActivityBatchTime: activityBatchTime,
	}
	subscriptionManager.Init()

	var stravaService = &api.StravaService{
		ClientId:         stravaClientId,
		SecretId:         stravaSecretId,
		StravaClient:     stravaClient,
		StravaRepository: stravaRepository,
	}

	var backgroundWorker = &BackgroundWorker{
		SubscriptionManager: subscriptionManager,
		StravaService:       stravaService,
		StravaRepository:    stravaRepository,
	}
	activitiesQueue := backgroundWorker.RunBackgroundWorker()

	router := mux.NewRouter()
	defaultApi := &api.DefaultApi{
		Router:          router,
		ApplicationSlug: applicationSlug,
	}

	var restApi = &api.SubscriptionApi{
		ActivitiesQueue: activitiesQueue,
		DefaultApi:      defaultApi,
	}
	restApi.RegisterHealth()
	restApi.RegisterApiRoutes()

	graphqlApi := graph.GraphqlApi{
		Resolvers: &graph.Resolver{
			ActivitiesQueue:     activitiesQueue,
			StravaService:       stravaService,
			SubscriptionManager: subscriptionManager,
			Repository:          stravaRepository,
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
			close(activitiesQueue)

			subscriptionManager.Close()
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
