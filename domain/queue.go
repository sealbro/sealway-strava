package domain

import "context"

type ActivitiesQueue struct {
	Channel chan StravaSubscriptionData
}

func (queue *ActivitiesQueue) Close(context.Context) error {
	close(queue.Channel)

	return nil
}
