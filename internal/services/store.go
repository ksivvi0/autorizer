package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Store struct {
	client *mongo.Client
}

func NewStoreInstance(uri string) (*Store, error) {
	clientOptions := options.Client().
		ApplyURI(uri)
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &Store{
		client: client,
	}, nil
}

func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *Store) WriteTokensInfo(token string) error {
	return errors.New("not implemented")
}
