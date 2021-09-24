package api

import (
	"encoding/json"
	"net/http"
	"sealway-strava/api/model"
)

type SubscriptionApi struct {
	Queue chan model.StravaSubscriptionData
	*DefaultApi
}

func (api *SubscriptionApi) RegisterApiRoutes() {
	var serverName = "api"
	api.Router.HandleFunc(api.Prefix(serverName, "/subscription"), api.subscription).Methods("POST")
}

// Methods

func (api *SubscriptionApi) subscription(w http.ResponseWriter, r *http.Request) {
	var data model.StravaSubscriptionData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	api.Queue <- data

	respondWithJSON(w, http.StatusCreated, nil)
}
