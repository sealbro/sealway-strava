package graph

import (
	"sealway-strava/domain"
	"sealway-strava/repository"
	"sealway-strava/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository          *repository.StravaRepository
	SubscriptionManager *usercase.SubscriptionManager
	StravaService       *usercase.StravaService
	ActivitiesQueue     chan domain.StravaSubscriptionData
}
