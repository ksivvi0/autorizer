package services

import (
	"authorizer/internal/helpers"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type StoreService interface {
	WriteTokensInfo(context.Context, tokenPair) ([]interface{}, error)
	GetTokensInfo(context.Context, string, string) ([]tokenPair, error)
	DropTokensInfo(context.Context, string, string) (int64, error)
}

type Store struct {
	client         *mongo.Client
	mainCollection *mongo.Collection
	cryptoKey      []byte
}

type refreshInfo struct {
	ID primitive.ObjectID `bson:"_id"`
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
	//if err = s.checkIndex("access_token_expired"); err != nil {
	//	return nil, err
	//}
	//if err = s.checkIndex("refresh_token_expired"); err != nil {
	//	return nil, err
	//}

	return s, nil
}

//func (s *Store) checkIndex(indexName string) error {
//	var indexes []bson.M
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//
//	cursor, err := s.mainCollection.Indexes().List(ctx)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err := cursor.Close(ctx); err != nil {
//			cancel()
//		}
//	}()
//
//	err = cursor.All(ctx, &indexes)
//	if err != nil {
//		return err
//	}
//	found := false
//	for _, v := range indexes {
//		if v["name"] == indexName {
//			found = true
//			break
//		}
//	}
//	if !found {
//		data := mongo.IndexModel{
//			Keys:    bson.M{indexName: 1},
//			Options: options.Index().SetExpireAfterSeconds(0),
//		}
//		_, err := s.mainCollection.Indexes().CreateOne(ctx, data)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *Store) WriteTokensInfo(ctx context.Context, pair tokenPair) ([]interface{}, error) {
	insertedIDs := make([]interface{}, 2)
	if pair.RefreshToken != "" {
		hashedTokenString, err := helpers.GetHash(pair.RefreshToken)
		if err != nil {
			return nil, err
		}
		refreshInfoResult, err := s.mainCollection.InsertOne(ctx, bson.M{
			"refresh_token_hash":    hashedTokenString,
			"refresh_token_expired": pair.RefreshExpired,
			"refresh_token_uid":     pair.RefreshUID,
		})
		if err != nil {
			return nil, err
		}
		insertedIDs[0] = refreshInfoResult.InsertedID
	}
	if pair.AccessUID != "" {
		accessInfoResult, err := s.mainCollection.InsertOne(ctx, bson.M{
			"access_token_uid":     pair.AccessUID,
			"access_token_expired": primitive.NewDateTimeFromTime(pair.AccessExpired),
			"refresh_token_uid":    pair.RefreshUID,
		})
		if err != nil {
			return nil, err
		}
		insertedIDs[0] = accessInfoResult.InsertedID
	}

	return insertedIDs, nil
}

func (s *Store) GetTokensInfo(ctx context.Context, key, uid string) ([]tokenPair, error) {
	cursor, err := s.mainCollection.Find(ctx, bson.D{{key, uid}})
	if err != nil {
		return nil, err
	}

	if cursor.RemainingBatchLength() == 0 {
		return nil, errors.New("failed to get access token")
	}

	pairs := make([]tokenPair, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var pair tokenPair
		if err := cursor.Decode(&pair); err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (s *Store) DropTokensInfo(ctx context.Context, key, value string) (int64, error) {
	accessResult, err := s.mainCollection.DeleteMany(ctx, bson.M{key: value})
	if err != nil {
		return -1, err
	}
	return accessResult.DeletedCount, nil
}
