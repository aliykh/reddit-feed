//go:generate mockgen -source mongo.go -destination mock/mongo_mock.go -package mock
package db

import (
	"context"
	"github.com/aliykh/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, res interface{}, opts ...*options.FindOneOptions) error
	Find(ctx context.Context, filter interface{}, res interface{}, opts ...*options.FindOptions) error
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline, res interface{}, opts ...*options.AggregateOptions) error
}

type dbCollection struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     *log.Factory
}

func New(logger *log.Factory, client *mongo.Client, dbName string, collName string) Collection {
	return &dbCollection{
		client:     client,
		collection: client.Database(dbName).Collection(collName),
		logger:     logger,
	}
}

func (m *dbCollection) Find(ctx context.Context, filter interface{}, res interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	cur, err := m.collection.Find(ctx, filter, opts...)

	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	return cur.All(ctx, res)
}

func (m *dbCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	return m.collection.InsertOne(ctx, document, opts...)
}

func (m *dbCollection) FindOne(ctx context.Context, filter interface{}, res interface{}, opts ...*options.FindOneOptions) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	return m.collection.FindOne(ctx, filter, opts...).Decode(res)
}
func (m *dbCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	return m.collection.CountDocuments(ctx, filter, opts...)
}
func (m *dbCollection) Aggregate(ctx context.Context, pipeline mongo.Pipeline, res interface{}, opts ...*options.AggregateOptions) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	cur, err := m.collection.Aggregate(ctx, pipeline, opts...)

	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	return cur.All(ctx, res)
}
