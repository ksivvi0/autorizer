package services

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type Store struct {
	client         *mongo.Client
	mainCollection *mongo.Collection
}

func NewStoreInstance(uri string) (*Store, error) {
	clientOptions := options.Client().
		ApplyURI(uri)
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	return &Store{
		client:         client,
		mainCollection: collection,
	}, nil
}

func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *Store) WriteTokensInfo(ctx context.Context, pair *tokenPair) (interface{}, error) {
	result, err := s.mainCollection.InsertOne(ctx, pair)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}
