package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"sealway-strava/domain"
	"sealway-strava/domain/strava"
	"sealway-strava/interfaces/graph/generated"
	"sealway-strava/interfaces/graph/model"
	"sealway-strava/pkg/logger"
)

func (r *mutationResolver) AddToken(ctx context.Context, tokens []*model.NewAthleteToken) (*string, error) {
	var errResult error

	for _, input := range tokens {
		err := r.Repository.UpsertToken(ctx, &domain.StravaToken{
			AthleteID: input.AthleteID,
			Refresh:   input.Refresh,
		})

		logger.Infof("Refresh token for [%d]", input.AthleteID)

		if err != nil {
			errResult = err
		}
	}

	return nil, errResult
}

func (r *mutationResolver) ResendSavedActivities(ctx context.Context, athleteIds []int64, limit int64) (*string, error) {
	activities, err := r.Repository.GetActivities(ctx, athleteIds, limit)
	if err == nil {
		r.SubscriptionManager.Notify(activities...)
	}

	return nil, err
}

func (r *mutationResolver) ReloadAthleteActivities(ctx context.Context, athleteIds []int64, before *int64, after *int64, page *int64, limit int64) (*string, error) {
	var errResult error

	for _, athleteId := range athleteIds {
		activities, err := r.StravaService.GetActivitiesByAthleteId(ctx, athleteId, before, after, page, limit)
		if err != nil {
			errResult = err
		} else {
			for _, activity := range activities {
				r.ActivitiesQueue <- domain.StravaSubscriptionData{
					AthleteId:  athleteId,
					ActivityId: activity.Id,
					Type:       "activity",
					Operation:  "create",
				}
			}
		}
	}

	return nil, errResult
}

func (r *queryResolver) Activity(ctx context.Context, id int64) (*strava.DetailedActivity, error) {
	return r.Repository.GetActivity(ctx, id)
}

func (r *queryResolver) Activities(ctx context.Context, athleteIds []int64, limit int64) ([]*strava.DetailedActivity, error) {
	return r.Repository.GetActivities(ctx, athleteIds, limit)
}

func (r *queryResolver) Token(ctx context.Context, athleteID int64) ([]*model.AthleteToken, error) {
	token, err := r.StravaService.RefreshToken(athleteID)
	if err != nil {
		return nil, err
	}

	return []*model.AthleteToken{{
		AthleteID: athleteID,
		Refresh:   *token,
	}}, nil
}

func (r *subscriptionResolver) Activities(ctx context.Context) (<-chan []*strava.DetailedActivity, error) {
	return r.SubscriptionManager.AddSubscriber(ctx)
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
