package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"net/http"
	"sealway-strava/strava"
)

type StravaService struct {
	StravaClient     *strava.APIClient
	ClientId         string
	SecretId         string
	StravaRepository *StravaRepository
}

func (service *StravaService) GetActivityById(token string, activityId int64) (*strava.DetailedActivity, error) {
	auth := context.WithValue(context.Background(), strava.ContextAccessToken, token)

	activity, _, err := service.StravaClient.ActivitiesApi.GetActivityById(auth, activityId, &strava.ActivitiesApiGetActivityByIdOpts{
		IncludeAllEfforts: optional.NewBool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("can't get activity %d : %w", activityId, err)
	}

	return &activity, nil
}

// TODO redis cache
func (service *StravaService) RefreshToken(athleteId int64) (*string, error) {
	stravaToken, err := service.StravaRepository.GetToken(athleteId)
	if err != nil {
		return nil, err
	}

	accessToken, err := service.refreshToken(stravaToken.Refresh)
	if err != nil {
		return nil, err
	}

	return accessToken, err
}

func (service *StravaService) refreshToken(refreshToken string) (*string, error) {
	url := "https://www.strava.com/api/v3/oauth/token"

	values := map[string]string{
		"client_id":     service.ClientId,
		"client_secret": service.SecretId,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	accessToken := res["access_token"].(string)

	return &accessToken, nil
}
