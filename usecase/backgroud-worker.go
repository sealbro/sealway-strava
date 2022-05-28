package usercase

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"sealway-strava/domain"
	"sealway-strava/domain/strava"
	"sealway-strava/pkg/closer"
	"sealway-strava/pkg/logger"
	"sealway-strava/repository"
	"time"
)

type BackgroundWorker struct {
	stravaRepository    *repository.StravaRepository
	stravaService       *StravaService
	subscriptionManager *SubscriptionManager
	queue               *domain.ActivitiesQueue
}

func MakeBackgroundWorker(collection *closer.CloserCollection, repository *repository.StravaRepository, service *StravaService, manager *SubscriptionManager) *domain.ActivitiesQueue {
	worker := &BackgroundWorker{
		stravaRepository:    repository,
		stravaService:       service,
		subscriptionManager: manager,
		queue: &domain.ActivitiesQueue{
			Channel: make(chan domain.StravaSubscriptionData),
		},
	}

	go worker.process(context.Background())

	collection.Add(worker.queue)

	return worker.queue
}

func (worker *BackgroundWorker) process(ctx context.Context) {
	for {
		data, ok := <-worker.queue.Channel
		if ok == false {
			break
		}

		go worker.processTaskWithRetry(ctx, data)
	}
}

func (worker *BackgroundWorker) processTaskWithRetry(ctx context.Context, data domain.StravaSubscriptionData) {
	err := retry.Do(
		func() error {
			logger.Infof("Start attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)
			err := worker.processTask(ctx, data)
			logger.Infof("Finish attempt process for activity [%d] with athlete [%d]", data.ActivityId, data.AthleteId)

			return err
		},
		retry.Context(ctx),
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

func (worker *BackgroundWorker) processTask(ctx context.Context, data domain.StravaSubscriptionData) error {
	athleteId := data.AthleteId

	if err := worker.stravaRepository.AddNewSubscription(ctx, &domain.StravaSubscription{
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		Data:     data,
	}); err != nil {
		return fmt.Errorf("BackgroundWorker - processTask - can't insert subscription for activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
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

			if activity, err = worker.stravaRepository.UpdateDetailedActivity(ctx, data.ActivityId, props); err != nil {
				return fmt.Errorf("BackgroundWorker - processTask - can't update activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		case "create":
			fallthrough
		default:
			if activity, err = worker.stravaService.SaveActivityById(ctx, athleteId, data.ActivityId); err != nil {
				return fmt.Errorf("BackgroundWorker - processTask - can't save activity [%d] for athlete [%d] : %s", data.ActivityId, athleteId, err.Error())
			}
		}
	}

	if activity != nil {
		worker.subscriptionManager.Notify(activity)
	}

	return nil
}
