package usercase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sealway-strava/domain/strava"
	"sealway-strava/pkg/batching"
	"sealway-strava/pkg/closer"
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

func MakeSubscriptionManager(collection *closer.CloserCollection, config *BatchConfig) *SubscriptionManager {
	manager := &SubscriptionManager{
		BatchConfig:   config,
		subscribers:   map[string]chan []*strava.DetailedActivity{},
		inputActivity: make(chan *strava.DetailedActivity),
		closed:        false,
	}

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

	collection.Add(manager)

	return manager
}

func (manager *SubscriptionManager) Notify(activities ...*strava.DetailedActivity) {
	if manager.closed {
		return
	}

	for _, activity := range activities {
		manager.inputActivity <- activity
	}
}

func (manager *SubscriptionManager) AddSubscriber(ctx context.Context) (chan []*strava.DetailedActivity, error) {
	if manager.closed {
		return nil, fmt.Errorf("SubscriptionManager closed the connection")
	}

	key := uuid.New().String()
	ch := make(chan []*strava.DetailedActivity)
	manager.subscribers[key] = ch

	logger.Infof("SubscriptionManager - Added new subscriber %s", key)

	go func() {
		<-ctx.Done()
		manager.RemoveSubscriber(key)
		logger.Infof("SubscriptionManager - Removed subscriber %s", key)
	}()

	return ch, nil
}

func (manager *SubscriptionManager) RemoveSubscriber(key string) {
	close(manager.subscribers[key])

	delete(manager.subscribers, key)
}

func (manager *SubscriptionManager) Close(context.Context) error {
	manager.closed = true

	for key := range manager.subscribers {
		manager.RemoveSubscriber(key)
	}

	close(manager.inputActivity)

	return nil
}
