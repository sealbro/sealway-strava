package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	api "sealway-strava/api/model"
	"sealway-strava/graph/generated"
	"sealway-strava/graph/model"
)

func (r *mutationResolver) AddToken(ctx context.Context, input model.NewAthleteToken) (int, error) {
	err := r.Repository.UpsertToken(&api.StravaToken{
		AthleteID: input.AthleteID,
		Access:    input.Access,
		Refresh:   input.Refresh,
		Expired:   input.Expired,
	})

	return 1, err
}

func (r *queryResolver) Activities(ctx context.Context) ([]*model.DetailedActivity, error) {
	activities, err := r.Repository.GetActivities()

	var result []*model.DetailedActivity

	for _, activity := range activities {
		result = append(result, &model.DetailedActivity{
			ID:          string(activity.Id),
			Distance:    float64(activity.Distance),
			Name:        activity.Name,
			Description: &activity.Description,
		})
	}

	return result, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
