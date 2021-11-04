package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"net/http"
	"sealway-strava/api/model"
	"sealway-strava/strava"
	"strconv"
	"strings"
)

type StravaService struct {
	StravaClient     *strava.APIClient
	ClientId         string
	SecretId         string
	StravaRepository *StravaRepository
}

func (service *StravaService) GetActivityById(athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	err := stravaQuota.CheckQuota()
	if err != nil {
		return nil, err
	}

	auth := service.GetStravaAuthContext(context.Background(), athleteId)

	activity, response, err := service.StravaClient.ActivitiesApi.GetActivityById(auth, activityId, &strava.ActivitiesApiGetActivityByIdOpts{
		IncludeAllEfforts: optional.NewBool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("can't get activity %d : %w", activityId, err)
	}

	updateQuota(response)

	return &activity, nil
}

func (service *StravaService) GetActivitiesByAthleteId(ctx context.Context, athleteId int64, before *int64, after *int64, page *int64, limit int64) ([]strava.SummaryActivity, error) {
	err := stravaQuota.CheckQuota()
	if err != nil {
		return nil, err
	}

	auth := service.GetStravaAuthContext(ctx, athleteId)

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

	activities, response, err := service.StravaClient.ActivitiesApi.GetLoggedInAthleteActivities(auth, &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
		Before:  beforeValue,
		After:   afterValue,
		Page:    pageValue,
		PerPage: perPageValue,
	})

	if err != nil {
		return nil, fmt.Errorf("can't get activities for %d : %w", athleteId, err)
	}

	updateQuota(response)

	return activities, nil
}

func (service *StravaService) GetStravaAuthContext(ctx context.Context, athleteId int64) context.Context {
	token, err := service.RefreshToken(athleteId)
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, strava.ContextAccessToken, *token)
}

// TODO redis cache
func (service *StravaService) RefreshToken(athleteId int64) (*string, error) {
	stravaToken, err := service.StravaRepository.GetToken(athleteId)
	if err != nil {
		return nil, err
	}

	accessToken, err := service.apiRefreshToken(stravaToken.Refresh)
	if err != nil {
		return nil, err
	}

	return accessToken, err
}

func (service *StravaService) apiRefreshToken(refreshToken string) (*string, error) {
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

func updateQuota(response *http.Response) {
	if response == nil {
		return
	}

	limitHeader := "X-Ratelimit-Limit"
	usageHeader := "X-Ratelimit-Usage"

	limits := response.Header[limitHeader]
	usages := response.Header[usageHeader]

	if len(limits) == 0 || len(usages) == 0 {
		return
	}

	limitValues := strings.Split(limits[0], ",")
	usageValues := strings.Split(usages[0], ",")

	limit15min, _ := strconv.Atoi(limitValues[0])
	limitDay, _ := strconv.Atoi(limitValues[1])
	usage15min, _ := strconv.Atoi(usageValues[0])
	usageDay, _ := strconv.Atoi(usageValues[1])

	stravaQuota = model.StravaQuota{
		Limit15min: limit15min,
		LimitDay:   limitDay,
		Usage15min: usage15min,
		UsageDay:   usageDay,
	}
}
