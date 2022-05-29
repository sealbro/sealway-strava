package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sealway-strava/domain"
	"sealway-strava/infrastructure"
	"sealway-strava/pkg/logger"
)

type SubscriptionApi struct {
	*DefaultApi

	ActivitiesQueue chan domain.StravaSubscriptionData
	StravaClient    *infrastructure.StravaClient
}

func MakeSubscriptionApi(queue *domain.ActivitiesQueue, api *DefaultApi, client *infrastructure.StravaClient) *SubscriptionApi {
	var restApi = &SubscriptionApi{
		ActivitiesQueue: queue.Channel,
		StravaClient:    client,
		DefaultApi:      api,
	}
	restApi.RegisterHealth()
	restApi.RegisterApiRoutes()

	return restApi
}

func (api *SubscriptionApi) RegisterApiRoutes() {
	urlPrefix := "api"
	api.Router.HandleFunc(api.Prefix(urlPrefix, "/quota"), api.quota).Methods("GET")
	api.Router.HandleFunc(api.Prefix(urlPrefix, "/subscription"), api.verify).Methods("GET")
	api.Router.HandleFunc(api.Prefix(urlPrefix, "/subscription"), api.subscription).Methods("POST")
}

func (api *SubscriptionApi) quota(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, &api.StravaClient.Quota)
}

// Methods
// https://developers.strava.com/docs/webhooks/

func (api *SubscriptionApi) verify(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["hub.challenge"]

	logger.Infof("Verify [%s]", r.URL.Path)

	respondWithJSON(w, http.StatusOK, &domain.StravaVerify{Challenge: keys[0]})
}

func (api *SubscriptionApi) subscription(w http.ResponseWriter, r *http.Request) {

	all, _ := io.ReadAll(r.Body)
	logger.Infof("Request to %s - %s", r.RequestURI, string(all))

	var data domain.StravaSubscriptionData
	reader := bytes.NewReader(all)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	api.ActivitiesQueue <- data

	respondWithJSON(w, http.StatusCreated, "Successful")
}
