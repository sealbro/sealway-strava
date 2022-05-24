package usercase

import (
	"fmt"
	"github.com/avast/retry-go"
	"sealway-strava/domain"
	"sealway-strava/domain/strava"
	"sealway-strava/pkg/logger"
	"sealway-strava/repository"
	"time"
)

type BackgroundWorker struct {
	StravaRepository    *repository.StravaRepository
	StravaService       *StravaService
	SubscriptionManager *SubscriptionManager
}

func MakeBackgroundWorker(repository *repository.StravaRepository, service *StravaService, manager *SubscriptionManager) *domain.ActivitiesQueue {
	worker := &BackgroundWorker{
		StravaRepository:    repository,
		StravaService:       service,
		SubscriptionManager: manager,
	}

	return worker.RunBackgroundWorker()
}

func (worker *BackgroundWorker) RunBackgroundWorker() *domain.ActivitiesQueue {
	inboundQueue := make(chan domain.StravaSubscriptionData)

	go worker.process(inboundQueue)

	return &domain.ActivitiesQueue{
		Channel: inboundQueue,
	}
}

func (worker *BackgroundWorker) process(inboundQueue chan domain.StravaSubscriptionData) {
	for {
		// check exit
		data, ok := <-inboundQueue
		if ok == false {
			break
		}

		err := retry.Do(
			func() error {
				logger.Infof("Start attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)
				activity, err := worker.processTask(data)
				if err != nil {
					logger.Error(err.Error())
				}

				if activity != nil {
					worker.SubscriptionManager.Notify([]*strava.DetailedActivity{activity})
				}

				logger.Infof("Finish attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)

				return err
			},
			retry.Attempts(5),
			retry.Delay(time.Minute),
			retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
				return time.Duration(n) * time.Minute
			}),
		)

		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func (worker *BackgroundWorker) processTask(data domain.StravaSubscriptionData) (*strava.DetailedActivity, error) {
	athleteId := data.AthleteId

	if err := worker.StravaRepository.AddNewSubscription(&domain.StravaSubscription{
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		Data:     data,
	}); err != nil {
		return nil, fmt.Errorf("BackgroundWorker - processTask - can't insert subscription for activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
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
				return nil, fmt.Errorf("BackgroundWorker - processTask - can't update activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		case "create":
			fallthrough
		default:
			if activity, err = worker.SaveActivityById(athleteId, data.ActivityId); err != nil {
				return nil, fmt.Errorf("BackgroundWorker - processTask - can't save activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		}
	}

	return activity, nil
}

func (worker *BackgroundWorker) SaveActivityById(athleteId int64, activityId int64) (*strava.DetailedActivity, error) {
	activity, err := worker.StravaService.GetActivityById(athleteId, activityId)
	if err != nil {
		return nil, fmt.Errorf("BackgroundWorker - SaveActivityById - GetActivityById: %s", err.Error())
	}

	err = worker.StravaRepository.AddDetailedActivity(activity)
	if err != nil {
		return nil, fmt.Errorf("BackgroundWorker - SaveActivityById - AddDetailedActivity: %s", err.Error())
	}

	return activity, nil
}
