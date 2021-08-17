package services

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
