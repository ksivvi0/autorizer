package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type StoreService interface {
	WriteTokensInfo(context.Context, *tokenPair) (interface{}, error)
	GetTokensInfo(context.Context, string) (*tokenPair, error)
}

type Store struct {
	client         *mongo.Client
	mainCollection *mongo.Collection
	cryptoKey      []byte
}

func NewStoreInstance(uri string) (*Store, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	s := &Store{
		client:         client,
		mainCollection: collection,
		cryptoKey:      []byte(os.Getenv("CRYPTO_SECRET")),
	}
	if err = s.checkIndex(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) checkIndex() error {
	var indexes []bson.M
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	cursor, err := s.mainCollection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			cancel()
		}
	}()

	err = cursor.All(ctx, &indexes)
	if err != nil {
		return err
	}
	found := false
	for _, v := range indexes {
		if v["name"] == "expires_at_1" {
			found = true
			break
		}
	}
	if !found {
		var iName string = "expires_at_1"
		var expiresAt int32 = 0
		opts := options.Index()
		opts.ExpireAfterSeconds = &expiresAt
		opts.Name = &iName

		tmp := mongo.IndexModel{
			Keys:    bson.M{"expires_at": 1},
			Options: opts,
		}
		_, err := s.mainCollection.Indexes().CreateOne(ctx, tmp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *Store) WriteTokensInfo(ctx context.Context, pair *tokenPair) (interface{}, error) {
	result, err := s.mainCollection.InsertOne(ctx, bson.D{
		{"access_uid", pair.AccessUID},
		{"expires_at", pair.AccessExpired.Format(time.UnixDate)},
	})
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (s *Store) GetTokensInfo(ctx context.Context, uid string) (*tokenPair, error) {
	pair := new(tokenPair)
	err := s.mainCollection.FindOne(ctx, bson.D{{"accessuid", uid}}).Decode(pair)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *Store) DropTokensInfo(ctx context.Context, uid string) (int64, error) {
	result, err := s.mainCollection.DeleteOne(ctx, bson.D{{"accessuid", uid}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
