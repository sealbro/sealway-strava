package domain

import (
	"time"
)

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
