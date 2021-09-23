package main

import (
	"strconv"
	"time"
)

type BackgroundWorker struct {
	stravaRepository *StravaRepository
	stravaService    *StravaService
}

func (worker *BackgroundWorker) RunBackgroundWorker() chan StravaSubscriptionData {
	stravaSubChannel := make(chan StravaSubscriptionData)

	go worker.process(stravaSubChannel)

	return stravaSubChannel
}

func (worker *BackgroundWorker) process(ch chan StravaSubscriptionData) {
	for {
		data, ok := <-ch
		if ok == false {
			break
		} else {
			athleteId, err := strconv.ParseInt(data.AthleteId, 10, 64)
			if err != nil {
				log.Errorf("can't convert %s to int64", data.AthleteId)
			} else {
				if err := worker.stravaRepository.AddNewSubscription(&StravaSubscription{
					CreateDate: time.Now(),
					Data:       data,
				}); err != nil {
					log.Error("can't insert subscription")
				}

				if err := worker.SaveActivityById(athleteId, data.ActivityId); err != nil {
					log.Errorf("can't save activity %d for %d : %w", data.ActivityId, athleteId, err)
				}
			}
		}
	}
}

func (worker *BackgroundWorker) SaveActivityById(athleteId int64, activityId int64) error {

	activity, err := worker.stravaService.GetActivityById(athleteId, activityId)
	if err != nil {
		return err
	}

	err = worker.stravaRepository.AddDetailedActivity(activity)
	if err != nil {
		return err
	}

	return nil
}
