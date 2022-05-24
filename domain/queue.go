package domain

type ActivitiesQueue struct {
	Channel chan StravaSubscriptionData
}

func (queue *ActivitiesQueue) Close() {
	close(queue.Channel)
}
