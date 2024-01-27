package datastore

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDataStore struct {
	client *mongo.Client
	dbName string
}

var (
	mongoDBInstance    *MongoDBDataStore
	mongoDBInstanceErr error
	mongoDBOnce        sync.Once
)

func GetMongoDBDataStoreInstance(uri, dbName string) (*MongoDBDataStore, error) {
	mongoDBOnce.Do(func() {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			mongoDBInstanceErr = err
			return
		}
		mongoDBInstance = &MongoDBDataStore{client: client, dbName: dbName}
	})
	return mongoDBInstance, mongoDBInstanceErr
}

func (m *MongoDBDataStore) Create(ctx context.Context, collection string, data interface{}) error {
	_, err := m.client.Database(m.dbName).Collection(collection).InsertOne(ctx, data)
	return err
}

func (m *MongoDBDataStore) Read(ctx context.Context, collection string, filter interface{}, result interface{}) error {
	return m.client.Database(m.dbName).Collection(collection).FindOne(ctx, filter).Decode(result)
}

func (m *MongoDBDataStore) Delete(ctx context.Context, collection string, filter interface{}) error {
	_, err := m.client.Database(m.dbName).Collection(collection).DeleteOne(ctx, filter)
	return err
}
