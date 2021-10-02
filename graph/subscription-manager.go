package graph

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sealway-strava/infra"
	"sealway-strava/strava"
)

type SubscriptionManager struct {
	subscribers map[string]chan []*strava.DetailedActivity
}

func (manager *SubscriptionManager) Init() {
	manager.subscribers = map[string]chan []*strava.DetailedActivity{}
}

func (manager *SubscriptionManager) Notify(activities []*strava.DetailedActivity) {
	var activityIds string
	for _, a := range activities {
		activityIds = fmt.Sprintf("%d,%s", a.ID, activityIds)
	}

	infra.Log.Infof("Send activities [%s] to subscribers [%d]", activityIds, len(manager.subscribers))
	for _, ch := range manager.subscribers {
		ch <- activities
	}
}

func (manager *SubscriptionManager) AddSubscriber(ctx context.Context) chan []*strava.DetailedActivity {
	key := uuid.New().String()
	ch := make(chan []*strava.DetailedActivity)
	manager.subscribers[key] = ch

	infra.Log.Infof("Added new subscriber %s", key)

	go func() {
		<-ctx.Done()
		manager.RemoveSubscriber(key)
		infra.Log.Infof("Removed subscriber %s", key)
	}()

	return ch
}

func (manager *SubscriptionManager) RemoveSubscriber(key string) {
	close(manager.subscribers[key])

	delete(manager.subscribers, key)
}
