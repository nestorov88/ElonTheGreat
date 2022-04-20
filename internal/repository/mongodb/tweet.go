package mongodb

import (
	"context"
	"fmt"
	"github.com/nestorov88/ElonTheGreat/pkg/tweet"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

//mongoRepository implements tweet.TweetRepository
type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

//newMongoClient is creating the client and connection to give mongodb uri
func newMongoClient(mongoURL string, mongoTiemout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTiemout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		fmt.Println(os.Getenv("MONGO_URI"))
		return nil, errors.Wrap(err, "repository.newMongoClient")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "repository.newMongoClient")
	}

	return client, err
}

// NewMongoRepository return instance of tweet.TweetRepository
func NewMongoRepository(mongoURL string, mongoDB string, mongoTimeout int) (tweet.TweetRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}

	repo.client = client
	return repo, nil
}

// Store is inserting new document to db
func (m *mongoRepository) Store(tweet *tweet.Tweet) error {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("tweets")

	_, err := collection.InsertOne(ctx, tweet)

	if err != nil {
		return errors.Wrap(err, "repository.Tweet.Store")
	}

	return nil
}

// StoreOrUpdate is inserting new tweet document to db or is updating if there is existing one
func (m *mongoRepository) StoreOrUpdate(t *tweet.Tweet) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("tweets")

	var res tweet.Tweet

	if err := collection.FindOne(ctx, bson.M{"id": t.ID}).Decode(&res); err != nil {

		return m.Store(t)

	} else {
		_, err = collection.UpdateOne(ctx, bson.M{"id": res.ID}, bson.D{{"$set", &t}})

		return err
	}

	return nil
}

// ByTimeOfTheDay query tweets collection by time of the day they are made and return their count
func (m *mongoRepository) ByTimeOfTheDay() ([]tweet.ByTimeOfTheDayResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("tweets")

	query := []bson.M{
		{
			"$project": bson.M{
				"hours": bson.M{
					"$hour": bson.M{
						"$dateFromString": bson.M{
							"dateString": bson.M{
								"$concat": []string{"$date", " ", "$time"},
							},
							"format": "%Y-%m-%d %H:%M:%S",
						},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$hours",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	var result []tweet.ByTimeOfTheDayResult

	c, err := collection.Aggregate(ctx, query)

	if err != nil {
		return result, errors.Wrap(err, "repository.MongoRepository.ByTimeOfTheDay")
	}

	err = c.All(ctx, &result)

	return result, errors.Wrap(err, "repository.MongoRepository.ByTimeOfTheDay")
}

// StatsByYear query tweets collection and return count of tweets, likes and retweets by year
func (m *mongoRepository) StatsByYear() ([]tweet.StatsByYearResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("tweets")

	query := []bson.M{
		{
			"$match": bson.M{
				"date": bson.M{
					"$exists": true,
					"$ne":     nil,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$year": bson.M{
						"$dateFromString": bson.M{
							"dateString": "$date",
							"format":     "%Y-%m-%d",
						},
					},
				},
				"sum":      bson.M{"$sum": 1},
				"likes":    bson.M{"$sum": "$likescount"},
				"retweets": bson.M{"$sum": "$retweetscount"},
			},
		},
		{
			"$project": bson.M{
				"date":     "$_id",
				"sum":      1,
				"likes":    1,
				"retweets": 1,
			},
		},
		{
			"$sort": bson.M{
				"date": 1,
			},
		},
	}

	var result []tweet.StatsByYearResult

	c, err := collection.Aggregate(ctx, query)

	if err != nil {

		return nil, errors.Wrap(err, "repository.MongoRepository.ByTimeOfTheDay")
	}

	c.All(ctx, &result)

	return result, errors.Wrap(err, "repository.MongoRepository.ByTimeOfTheDay")
}
