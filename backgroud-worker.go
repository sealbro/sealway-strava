package main

import (
	"fmt"
	"github.com/avast/retry-go"
	"sealway-strava/api"
	"sealway-strava/api/model"
	"sealway-strava/infra"
	"sealway-strava/strava"
	"strconv"
	"time"
)

type BackgroundWorker struct {
	StravaRepository *api.StravaRepository
	StravaService    *api.StravaService
}

func (worker *BackgroundWorker) RunBackgroundWorker() (chan model.StravaSubscriptionData, chan []*strava.DetailedActivity) {
	inboundQueue := make(chan model.StravaSubscriptionData)
	outboundQueue := make(chan []*strava.DetailedActivity)

	go worker.process(inboundQueue, outboundQueue)

	return inboundQueue, outboundQueue
}

func (worker *BackgroundWorker) process(inboundQueue chan model.StravaSubscriptionData, outboundQueue chan []*strava.DetailedActivity) {
	for {
		// check exit
		data, ok := <-inboundQueue
		if ok == false {
			break
		}

		err := retry.Do(
			func() error {
				infra.Log.Infof("Start worker for activity [%d] for athlete [%s]", data.ActivityId, data.AthleteId)
				activity, err := worker.processTask(data)
				if err != nil {
					infra.Log.Error(err.Error())
				}

				if activity != nil {
					outboundQueue <- []*strava.DetailedActivity{activity}
				}

				infra.Log.Infof("Finish worker for activity [%d] for athlete [%s]", data.ActivityId, data.AthleteId)

				return err
			},
			retry.Attempts(5),
			retry.Delay(time.Minute),
			retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
				return time.Duration(n) * time.Minute
			}),
		)

		if err != nil {
			infra.Log.Error(err.Error())
		}
	}
}

func (worker *BackgroundWorker) processTask(data model.StravaSubscriptionData) (*strava.DetailedActivity, error) {
	// convert athlete id
	athleteId, err := strconv.ParseInt(data.AthleteId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't convert %s to int64", data.AthleteId)
	}

	// save subscription
	if err := worker.StravaRepository.AddNewSubscription(&model.StravaSubscription{
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		Data:     data,
	}); err != nil {
		return nil, fmt.Errorf("can't insert subscription for activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
	}

	var activity *strava.DetailedActivity

	if data.Type == "activity" {
		switch data.Operation {
		case "update":
			props := map[string]interface{}{}
			for key, value := range data.Updates {
				dbPropName := key
				switch key {
				case "title":
					dbPropName = "name"
				}

				props[dbPropName] = value
			}

			if activity, err = worker.StravaRepository.UpdateDetailedActivity(data.ActivityId, props); err != nil {
				return nil, fmt.Errorf("can't update activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		case "create":
			fallthrough
		default:
			if activity, err = worker.SaveActivityById(athleteId, data.ActivityId); err != nil {
				return nil, fmt.Errorf("can't save activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		}
	}

	return activity, nil
}

func (worker *BackgroundWorker) SaveActivityById(athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	// TODO redis cache
	stravaToken, err := worker.StravaRepository.GetToken(athleteId)
	if err != nil {
		return nil, err
	}

	accessToken, err := worker.StravaService.RefreshToken(stravaToken.Refresh)
	if err != nil {
		return nil, err
	}

	activity, err := worker.StravaService.GetActivityById(*accessToken, activityId)
	if err != nil {
		return nil, err
	}

	err = worker.StravaRepository.AddDetailedActivity(activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
