package domain

import (
	"fmt"
	"time"
)

type StravaQuota struct {
	Limit15min int       `json:"limit_15_min"`
	LimitDay   int       `json:"limit_day"`
	Usage15min int       `json:"usage_15_min"`
	UsageDay   int       `json:"usage_day"`
	LastUpdate time.Time `json:"last_update"`
}

type StravaVerify struct {
	Challenge string `json:"hub.challenge"`
}

type StravaSubscription struct {
	ExpireAt time.Time              `bson:"expire_at" json:"expire_at"`
	Data     StravaSubscriptionData `bson:"data" json:"data"`
}

type StravaSubscriptionData struct {
	ActivityId     int64             `bson:"object_id" json:"object_id"`
	Type           string            `bson:"object_type" json:"object_type"`
	Operation      string            `bson:"aspect_type" json:"aspect_type"`
	AthleteId      int64             `bson:"owner_id" json:"owner_id"`
	Updates        map[string]string `bson:"updates" json:"updates"`
	SubscriptionId int64             `bson:"subscription_id" json:"subscription_id"`
	EventTime      int64             `bson:"event_time" json:"event_time"`
}

type StravaToken struct {
	AthleteID int64  `bson:"_id" json:"athlete_id"`
	Refresh   string `bson:"refresh" json:"refresh"`
}

func (quota *StravaQuota) CheckQuota() error {
	diff := time.Since(quota.LastUpdate)

	if quota.Limit15min-quota.Usage15min <= 0 && diff.Minutes() < 15 {
		return fmt.Errorf("15min quota [%d/%d]", quota.Usage15min, quota.Limit15min)
	}

	if quota.LimitDay-quota.UsageDay <= 0 && diff.Hours() < 24 {
		return fmt.Errorf("day quota [%d/%d]", quota.UsageDay, quota.LimitDay)
	}

	return nil
}
