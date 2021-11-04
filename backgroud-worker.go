package main

import (
	"fmt"
	"github.com/avast/retry-go"
	"sealway-strava/api"
	"sealway-strava/api/model"
	"sealway-strava/graph"
	"sealway-strava/infra"
	"sealway-strava/strava"
	"time"
)

type BackgroundWorker struct {
	StravaRepository    *api.StravaRepository
	StravaService       *api.StravaService
	SubscriptionManager *graph.SubscriptionManager
}

func (worker *BackgroundWorker) RunBackgroundWorker() chan model.StravaSubscriptionData {
	inboundQueue := make(chan model.StravaSubscriptionData)

	go worker.process(inboundQueue)

	return inboundQueue
}

func (worker *BackgroundWorker) process(inboundQueue chan model.StravaSubscriptionData) {
	for {
		// check exit
		data, ok := <-inboundQueue
		if ok == false {
			break
		}

		err := retry.Do(
			func() error {
				infra.Log.Infof("Start attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)
				activity, err := worker.processTask(data)
				if err != nil {
					infra.Log.Error(err.Error())
				}

				if activity != nil {
					worker.SubscriptionManager.Notify([]*strava.DetailedActivity{activity})
				}

				infra.Log.Infof("Finish attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)

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
	athleteId := data.AthleteId

	// save subscription
	if err := worker.StravaRepository.AddNewSubscription(&model.StravaSubscription{
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		Data:     data,
	}); err != nil {
		return nil, fmt.Errorf("can't insert subscription for activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
	}

	var activity *strava.DetailedActivity
	var err error

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
	activity, err := worker.StravaService.GetActivityById(athleteId, activityId)
	if err != nil {
		return nil, err
	}

	err = worker.StravaRepository.AddDetailedActivity(activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
