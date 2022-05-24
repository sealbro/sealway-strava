package infrastructure

import "sealway-strava/domain/strava"

type StravaClient struct {
	*strava.APIClient
}

func MakeStravaClient() *StravaClient {
	apiClient := strava.NewAPIClient(strava.NewConfiguration())

	return &StravaClient{
		APIClient: apiClient,
	}
}
