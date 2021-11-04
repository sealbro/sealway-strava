package graph

import (
	"sealway-strava/api"
	"sealway-strava/api/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository          *api.StravaRepository
	SubscriptionManager *SubscriptionManager
	StravaService       *api.StravaService
	ActivitiesQueue     chan model.StravaSubscriptionData
}
