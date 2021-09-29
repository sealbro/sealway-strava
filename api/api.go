package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sealway-strava/infra"
)

type DefaultApi struct {
	Router          *mux.Router
	ApplicationSlug string
}

func (api *SubscriptionApi) RegisterHealth() {
	api.Router.HandleFunc("/health", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	}).Methods("GET")
}

func (api *DefaultApi) Prefix(serverName string, path string) string {
	return fmt.Sprintf("/%s/%s%s", api.ApplicationSlug, serverName, path)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	infra.Log.Error(message)

	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
