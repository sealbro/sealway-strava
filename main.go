package main

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"sealway-strava/strava"
	"time"
)

var applicationSlug = "integration-strava"
var log = &Logger{}

// ENVs
var stravaClientId = os.Getenv("SEALWAY_Services__Strava__Client")
var stravaSecretId = os.Getenv("SEALWAY_Services__Strava__Secret")
var connectionString = os.Getenv("SEALWAY_ConnectionStrings__Mongo__Connection")
var port = ":8080"

func main() {
	stravaClient := strava.NewAPIClient(strava.NewConfiguration())

	ctx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	err, stravaRepository := InitStravaRepository(connectionString, ctx)
	if err != nil {
		panic(err)
	}

	var stravaService = &StravaService{
		stravaClient: stravaClient,
	}

	var backgroundWorker = &BackgroundWorker{
		stravaService:    stravaService,
		stravaRepository: stravaRepository,
	}
	queue := backgroundWorker.RunBackgroundWorker()

	router := mux.NewRouter()
	var api = &SubscriptionApi{
		queue:           queue,
		applicationSlug: applicationSlug,
		router:          router,
	}
	api.RegisterRoutes()

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	var graceful = &Graceful{
		startAction: func() error {
			return srv.ListenAndServe()
		},
		deferAction: func(ctx context.Context) error {
			cancelMongo()
			stravaRepository.client.Disconnect(ctx)
			close(queue)

			return nil
		},
		shutdownAction: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	}

	graceful.RunAndWait()
}
