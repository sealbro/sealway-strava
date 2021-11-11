package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sealway-strava/api/model"
	"sealway-strava/infra"
	"time"
)

var stravaQuota = model.StravaQuota{
	Limit15min: 100,
	LimitDay:   1000,
	Usage15min: 0,
	UsageDay:   0,
	LastUpdate: time.Now(),
}

type SubscriptionApi struct {
	ActivitiesQueue chan model.StravaSubscriptionData
	*DefaultApi
}

func (api *SubscriptionApi) RegisterApiRoutes() {
	var serverName = "api"
	api.Router.HandleFunc(api.Prefix(serverName, "/healthz"), api.subscription).Methods("GET")
	api.Router.HandleFunc(api.Prefix(serverName, "/quota"), api.quota).Methods("GET")
	api.Router.HandleFunc(api.Prefix(serverName, "/subscription"), api.verify).Methods("GET")
	api.Router.HandleFunc(api.Prefix(serverName, "/subscription"), api.subscription).Methods("POST")
}

// Methods
// https://developers.strava.com/docs/webhooks/

func (api *SubscriptionApi) verify(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["hub.challenge"]

	infra.Log.Infof("Verify [%s]", r.URL.Path)

	respondWithJSON(w, http.StatusOK, &model.StravaVerify{Challenge: keys[0]})
}

func (api *SubscriptionApi) subscription(w http.ResponseWriter, r *http.Request) {

	all, _ := io.ReadAll(r.Body)
	infra.Log.Infof("Request to %s - %s", r.RequestURI, string(all))

	var data model.StravaSubscriptionData
	reader := bytes.NewReader(all)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	api.ActivitiesQueue <- data

	respondWithJSON(w, http.StatusCreated, "successful")
}

func (api *SubscriptionApi) quota(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, &stravaQuota)
}
