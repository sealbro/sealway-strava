package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type SubscriptionApi struct {
	applicationSlug string
	router          *mux.Router
	queue           chan StravaSubscriptionData
}

type StravaSubscriptionData struct {
	ActivityId     int64             `bson:"object_id" json:"object_id"`
	Type           string            `bson:"object_type" json:"object_type"`
	Operation      string            `bson:"aspect_type" json:"aspect_type"`
	AthleteId      string            `bson:"owner_id" json:"owner_id"`
	Updates        map[string]string `bson:"updates" json:"updates"`
	SubscriptionId int64             `bson:"subscription_id" json:"subscription_id"`
	EventTime      int64             `bson:"event_time" json:"event_time"`
}

func (api *SubscriptionApi) RegisterRoutes() {

	api.router.HandleFunc(api.prefix("/subscription"), api.subscription).Methods("POST")
}

func (api *SubscriptionApi) prefix(path string) string {
	return fmt.Sprintf("/%s/api%s", api.applicationSlug, path)
}

// Methods

func (api *SubscriptionApi) subscription(w http.ResponseWriter, r *http.Request) {
	var data StravaSubscriptionData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	api.queue <- data

	respondWithJSON(w, http.StatusCreated, nil)
}
