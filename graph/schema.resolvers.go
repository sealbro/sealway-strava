package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	api "sealway-strava/api/model"
	"sealway-strava/graph/generated"
	"sealway-strava/graph/model"
	"sealway-strava/infra"
	"sealway-strava/strava"
)

func (r *mutationResolver) AddToken(ctx context.Context, input model.NewAthleteToken) (*string, error) {
	err := r.Repository.UpsertToken(ctx, &api.StravaToken{
		AthleteID: input.AthleteID,
		Refresh:   input.Refresh,
	})

	infra.Log.Infof("Refresh token for %s", input.AthleteID)

	return nil, err
}

func (r *mutationResolver) ResendSavedActivities(ctx context.Context, athleteIds []int64, limit int64) (*string, error) {
	activities, err := r.Repository.GetActivities(ctx, athleteIds, limit)
	if err == nil {
		r.SubscriptionManager.Notify(activities)
	}

	return nil, err
}

func (r *queryResolver) Activity(ctx context.Context, id int64) (*strava.DetailedActivity, error) {
	return r.Repository.GetActivity(ctx, id)
}

func (r *queryResolver) Activities(ctx context.Context, athleteIds []int64, limit int64) ([]*strava.DetailedActivity, error) {
	return r.Repository.GetActivities(ctx, athleteIds, limit)
}

func (r *queryResolver) Token(ctx context.Context, athleteID int64) (*model.AthleteToken, error) {
	token, err := r.StravaService.RefreshToken(athleteID)
	if err != nil {
		return nil, err
	}

	return &model.AthleteToken{
		Refresh: *token,
	}, nil
}

func (r *subscriptionResolver) Activities(ctx context.Context) (<-chan []*strava.DetailedActivity, error) {
	return r.SubscriptionManager.AddSubscriber(ctx), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
