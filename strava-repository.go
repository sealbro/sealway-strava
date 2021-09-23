package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sealway-strava/strava"
	"time"
)

var stravaDataBaseName = "Strava"
var stravaSubscriptionCollectionName = "Subscription"
var stravaActivityCollectionName = "Activity"

type StravaRepository struct {
	client *mongo.Client
}

type StravaSubscription struct {
	CreateDate time.Time              `bson:"create_date" json:"create_date"`
	Data       StravaSubscriptionData `bson:"data" json:"data"`
}

func InitStravaRepository(connectionString string, ctx context.Context) (error, *StravaRepository) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err.Error())
		return err, nil
	}

	return err, &StravaRepository{
		client: client,
	}
}

// Operations

func (repository *StravaRepository) AddNewSubscription(data *StravaSubscription) error {
	collection := repository.client.Database(stravaDataBaseName).Collection(stravaSubscriptionCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, data)

	return err
}

func (repository *StravaRepository) AddDetailedActivity(activity *strava.DetailedActivity) error {
	collection := repository.client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, activity)

	if err != nil {
		return fmt.Errorf("can't insert activity %d : %w", activity.Id, err)
	}

	return nil
}

// Helpers

func createTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
