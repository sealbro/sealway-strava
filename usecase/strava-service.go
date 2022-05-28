package usercase

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"sealway-strava/domain/strava"
	"sealway-strava/infrastructure"
	"sealway-strava/pkg/closer"
	"sealway-strava/repository"
)

type StravaService struct {
	client     *infrastructure.StravaClient
	repository *repository.StravaRepository
}

func MakeStravaService(collection *closer.CloserCollection, client *infrastructure.StravaClient, repository *repository.StravaRepository) *StravaService {
	service := &StravaService{
		client:     client,
		repository: repository,
	}

	collection.Add(service)

	return service
}

func (service *StravaService) Close(ctx context.Context) error {
	service.client.Close()

	return service.repository.Close(ctx)
}

func (service *StravaService) SaveActivityById(ctx context.Context, athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	activity, err := service.GetActivityById(ctx, athleteId, activityId)
	if err != nil {
		return nil, fmt.Errorf("BackgroundWorker - SaveActivityById - GetActivityById: %s", err.Error())
	}

	err = service.repository.AddDetailedActivity(ctx, activity)
	if err != nil {
		return nil, fmt.Errorf("BackgroundWorker - SaveActivityById - AddDetailedActivity: %s", err.Error())
	}

	return activity, nil
}

func (service *StravaService) GetActivityById(ctx context.Context, athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	err := service.client.CheckQuota()
	if err != nil {
		return nil, err
	}

	auth := service.getStravaAuthContext(ctx, athleteId)

	activity, response, err := service.client.ActivitiesApi.GetActivityById(auth, activityId, &strava.ActivitiesApiGetActivityByIdOpts{
		IncludeAllEfforts: optional.NewBool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("StravaService - GetActivityById - can't get activity %d : %w", activityId, err)
	}

	service.client.UpdateQuota(response)

	return &activity, nil
}

func (service *StravaService) GetActivitiesByAthleteId(ctx context.Context, athleteId int64, before *int64, after *int64, page *int64, limit int64) ([]strava.SummaryActivity, error) {
	err := service.client.CheckQuota()
	if err != nil {
		return nil, err
	}

	auth := service.getStravaAuthContext(ctx, athleteId)

	var beforeValue optional.Int32
	var afterValue optional.Int32
	var pageValue optional.Int32
	perPageValue := optional.NewInt32(int32(limit))

	if before != nil {
		beforeValue = optional.NewInt32(int32(*before))
	}
	if after != nil {
		afterValue = optional.NewInt32(int32(*after))
	}
	if page != nil {
		pageValue = optional.NewInt32(int32(*page))
	}

	activities, response, err := service.client.ActivitiesApi.GetLoggedInAthleteActivities(auth, &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
		Before:  beforeValue,
		After:   afterValue,
		Page:    pageValue,
		PerPage: perPageValue,
	})

	if err != nil {
		return nil, fmt.Errorf("StravaService - GetActivitiesByAthleteId - can't get activities for %d : %w", athleteId, err)
	}

	service.client.UpdateQuota(response)

	return activities, nil
}

func (service *StravaService) getStravaAuthContext(ctx context.Context, athleteId int64) context.Context {
	token, err := service.RefreshToken(ctx, athleteId)
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, strava.ContextAccessToken, *token)
}

// TODO redis cache
func (service *StravaService) RefreshToken(ctx context.Context, athleteId int64) (*string, error) {
	stravaToken, err := service.repository.GetToken(ctx, athleteId)
	if err != nil {
		return nil, err
	}

	accessToken, err := service.client.RefreshToken(ctx, stravaToken.Refresh)
	if err != nil {
		return nil, err
	}

	return accessToken, err
}
