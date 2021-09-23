package main

import (
	"context"
	"fmt"
	"sealway-strava/strava"
)

type StravaService struct {
	stravaClient *strava.APIClient
}

func (service *StravaService) GetActivityById(athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	// TODO получать токен по AthleteId

	token := "138c63016d4fad5bf673f427bd9a52d972408216"
	auth := context.WithValue(context.Background(), strava.ContextAccessToken, token)

	activity, _, err := service.stravaClient.ActivitiesApi.GetActivityById(auth, activityId, nil)

	if err != nil {
		return nil, fmt.Errorf("can't get activity %d : %w", activityId, err)
	}

	return &activity, nil
}
