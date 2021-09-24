package main

import (
	"sealway-strava/api"
	"sealway-strava/api/model"
	"sealway-strava/infra"
	"strconv"
	"time"
)

type BackgroundWorker struct {
	StravaRepository *api.StravaRepository
	StravaService    *api.StravaService
}

func (worker *BackgroundWorker) RunBackgroundWorker() chan model.StravaSubscriptionData {
	stravaSubChannel := make(chan model.StravaSubscriptionData)

	go worker.process(stravaSubChannel)

	return stravaSubChannel
}

func (worker *BackgroundWorker) process(ch chan model.StravaSubscriptionData) {
	for {
		// TODO отдельную очередь для update/create c retry и батчем по времени
		data, ok := <-ch
		if ok == false {
			break
		} else {
			// TODO to process if object_type = activity and aspect_type (create or update)
			athleteId, err := strconv.ParseInt(data.AthleteId, 10, 64)
			if err != nil {
				infra.Log.Errorf("can't convert %s to int64", data.AthleteId)
			} else {
				// TODO add expire index, add to DB if error before save
				if err := worker.StravaRepository.AddNewSubscription(&model.StravaSubscription{
					CreateDate: time.Now(),
					Data:       data,
				}); err != nil {
					infra.Log.Error("can't insert subscription")
				}

				// TODO save if create
				if err := worker.SaveActivityById(athleteId, data.ActivityId); err != nil {
					infra.Log.Errorf("can't save activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
				}

				// TODO update changed properties
			}
		}
	}
}

func (worker *BackgroundWorker) SaveActivityById(athleteId int64, activityId int64) error {
	stravaToken, err := worker.StravaRepository.GetToken(athleteId)
	if err != nil {
		// TODO то в этом случае нужно сохранить subscribe
		return err
	}

	accessToken, err := worker.StravaService.RefreshToken(stravaToken.Refresh)
	if err != nil {
		return err
	}

	// TODO Cache or Update in DB
	//worker.StravaRepository.UpsertToken(&api.StravaToken{
	//	AthleteID: athleteId,
	//	Access:   *accessToken,
	//	Refresh: stravaToken.Refresh,
	//	Expired:   "",
	//})

	activity, err := worker.StravaService.GetActivityById(*accessToken, activityId)
	if err != nil {
		return err
	}

	err = worker.StravaRepository.AddDetailedActivity(activity)
	if err != nil {
		return err
	}

	return nil
}
