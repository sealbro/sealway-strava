package usercase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sealway-strava/domain/strava"
	"sealway-strava/pkg/batching"
	"sealway-strava/pkg/logger"
	"time"
)

type BatchConfig struct {
	ActivityBatchSize int
	ActivityBatchTime time.Duration
}

type SubscriptionManager struct {
	*BatchConfig

	subscribers      map[string]chan []*strava.DetailedActivity
	outputActivities chan []*strava.DetailedActivity
	inputActivity    chan *strava.DetailedActivity
	closed           bool
}

func MakeSubscriptionManager(config *BatchConfig) *SubscriptionManager {
	manager := &SubscriptionManager{
		BatchConfig: config,
	}

	manager.initialize()

	return manager
}

func (manager *SubscriptionManager) initialize() {
	manager.subscribers = map[string]chan []*strava.DetailedActivity{}
	manager.inputActivity = make(chan *strava.DetailedActivity)
	manager.outputActivities = batching.Process(manager.inputActivity, manager.ActivityBatchSize, manager.ActivityBatchTime)

	go func() {
		for activities := range manager.outputActivities {
			var activityIds string
			for _, a := range activities {
				activityIds = fmt.Sprintf("%d,%s", a.ID, activityIds)
			}

			logger.Infof("Send activities [%s] to subscribers [%d]", activityIds, len(manager.subscribers))
			for _, subscriber := range manager.subscribers {
				subscriber <- activities
			}
		}
	}()
}

func (manager *SubscriptionManager) Notify(activities ...*strava.DetailedActivity) {
	if manager.closed {
		logger.Fatal("inputActivity was closed")
		return
	}

	for _, activity := range activities {
		manager.inputActivity <- activity
	}
}

func (manager *SubscriptionManager) AddSubscriber(ctx context.Context) chan []*strava.DetailedActivity {
	key := uuid.New().String()
	ch := make(chan []*strava.DetailedActivity)
	manager.subscribers[key] = ch

	logger.Infof("Added new subscriber %s", key)

	go func() {
		<-ctx.Done()
		manager.RemoveSubscriber(key)
		logger.Infof("Removed subscriber %s", key)
	}()

	return ch
}

func (manager *SubscriptionManager) RemoveSubscriber(key string) {
	close(manager.subscribers[key])

	delete(manager.subscribers, key)
}

func (manager *SubscriptionManager) Close() {
	manager.closed = true

	close(manager.inputActivity)
}
