package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sealway-strava/domain"
	"sealway-strava/domain/strava"
	"sealway-strava/pkg/logger"
	"time"
)

var stravaDataBaseName = "Strava"
var stravaSubscriptionCollectionName = "Subscription"
var stravaActivityCollectionName = "Activity"
var stravaTokenCollectionName = "Token"

type MongoConfig struct {
	ConnectionString string
}

type StravaRepository struct {
	*mongo.Client
}

func MakeStravaRepository(config *MongoConfig) (*StravaRepository, error) {
	ctx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelMongo()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		return nil, err
	}

	repository := &StravaRepository{
		Client: client,
	}

	repository.addIndex(ctx, stravaDataBaseName, stravaActivityCollectionName, bson.M{"athlete.id": 1}, nil)
	repository.addIndex(ctx, stravaDataBaseName, stravaSubscriptionCollectionName, bson.D{{"expire_at", 1}}, options.Index().SetExpireAfterSeconds(0))
	repository.addIndex(ctx, stravaDataBaseName, stravaSubscriptionCollectionName, bson.D{{"data.owner_id", 1}, {"data.object_id", 1}}, nil)

	return repository, err
}

func (repository *StravaRepository) Close(ctx context.Context) error {
	return repository.Client.Disconnect(ctx)
}

func (repository *StravaRepository) addIndex(ctx context.Context, dbName string, collection string, indexKeys interface{}, opt *options.IndexOptions) {
	serviceCollection := repository.Client.Database(dbName).Collection(collection)
	indexName, err := serviceCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    indexKeys,
		Options: opt,
	})
	if err != nil {
		logger.Fatalf("can't create index %s -> %s -> %v: %v", collection, dbName, indexKeys, err)
	}

	logger.Tracef("Index created: %s", indexName)
}

// Operations

func (repository *StravaRepository) AddNewSubscription(ctx context.Context, data *domain.StravaSubscription) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaSubscriptionCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()
	_, err := collection.InsertOne(ctx, data)

	return err
}

func (repository *StravaRepository) UpdateDetailedActivity(ctx context.Context, activityID int64, updates interface{}) (*strava.DetailedActivity, error) {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()
	upd, err := collection.UpdateByID(ctx, activityID, bson.M{"$set": updates})

	if err != nil {
		return nil, err
	}

	if upd.ModifiedCount > 0 {
		return repository.GetActivity(ctx, activityID)
	}

	return nil, nil
}

func (repository *StravaRepository) AddDetailedActivity(ctx context.Context, activity *strava.DetailedActivity) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()
	collection.DeleteOne(ctx, bson.D{{"_id", activity.ID}})

	_, err := collection.InsertOne(ctx, activity)

	return err
}

func (repository *StravaRepository) UpsertToken(ctx context.Context, token domain.StravaToken) error {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaTokenCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", token.AthleteID}}
	update := bson.D{{"$set", token}}

	_, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return fmt.Errorf("can't insert token %d : %w", token.AthleteID, err)
	}

	return nil
}

func (repository *StravaRepository) GetToken(ctx context.Context, athleteId int64) (domain.StravaToken, error) {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaTokenCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()

	var token domain.StravaToken
	err := collection.FindOne(ctx, bson.D{{"_id", athleteId}}).Decode(&token)

	return token, err
}

func (repository *StravaRepository) GetActivity(ctx context.Context, activityId int64) (*strava.DetailedActivity, error) {
	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(ctx)
	defer cancel()

	var activity *strava.DetailedActivity
	err := collection.FindOne(ctx, bson.D{{"_id", activityId}}).Decode(&activity)

	return activity, err
}

func (repository *StravaRepository) GetActivities(innerCtx context.Context, athleteIds []int64, limit int64) ([]*strava.DetailedActivity, error) {
	if len(athleteIds) == 0 {
		return nil, fmt.Errorf("stravaRepository - GetActivities - 'athleteIds' is empty")
	}

	collection := repository.Client.Database(stravaDataBaseName).Collection(stravaActivityCollectionName)
	ctx, cancel := createTimeoutFromInnerContext(innerCtx)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{"athlete.id": bson.M{"$in": athleteIds}})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	var activities []*strava.DetailedActivity
	for cursor.Next(ctx) && limit > 0 {
		limit--
		var activity *strava.DetailedActivity
		err := cursor.Decode(&activity)
		if err != nil {
			logger.Tracef("decode activity : %s", err.Error())
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func createTimeoutFromInnerContext(innerCtx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(innerCtx, 15*time.Second)
}
