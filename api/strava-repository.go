package api

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sealway-strava/api/model"
	"sealway-strava/infra"
	"sealway-strava/strava"
	"time"
)

var stravaDataBaseName = "Strava"
var stravaSubscriptionCollectionName = "Subscription"
var stravaActivityCollectionName = "Activity"
var stravaTokenCollectionName = "Token"

type StravaRepository struct {
	Client *mongo.Client
}

func InitStravaRepository(connectionString string, ctx context.Context) (error, *StravaRepository) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		infra.Log.Fatal(err.Error())
		return err, nil
	}

	return err, &StravaRepository{
		Client: client,
	}
}

// Operations

func (repository *StravaRepository) AddNewSubscription(data *model.StravaSubscription) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaSubscriptionCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, data)

	return err
}

func (repository *StravaRepository) AddDetailedActivity(activity *strava.DetailedActivity) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, activity)

	if err != nil {
		return fmt.Errorf("can't insert activity %d : %w", activity.Id, err)
	}

	return nil
}

func (repository *StravaRepository) UpsertToken(token *model.StravaToken) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaTokenCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, token)

	if err != nil {
		return fmt.Errorf("can't insert token %d : %w", token.AthleteID, err)
	}

	return nil
}

func (repository *StravaRepository) GetToken(athleteId int64) (*model.StravaToken, error) {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaTokenCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()

	var token *model.StravaToken
	err := collection.FindOne(ctx, bson.D{{"_id", athleteId}}).Decode(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (repository *StravaRepository) GetActivities() ([]*strava.DetailedActivity, error) {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutContext()
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var activities []*strava.DetailedActivity
	for cursor.Next(ctx) {
		var activity *strava.DetailedActivity
		err := cursor.Decode(&activity)
		if err != nil {
			infra.Log.Errorf("decode activity : %s", err.Error())
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// Helpers

func createTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
