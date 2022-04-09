package repository

import (
	"context"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	dbName         = "reddit-feed"
	collectionName = "posts"
)

type Repository interface {
	Create(context.Context, *models.Post) (*models.Post, error)
}

type repo struct {
	logger   *log.Factory
	dbClient *mongo.Client
}

func New(logger *log.Factory, dbClient *mongo.Client) *repo {
	return &repo{
		logger:   logger,
		dbClient: dbClient,
	}
}

func (r *repo) Create(ctx context.Context, m *models.Post) (*models.Post, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := r.getCollection()

	res, err := coll.InsertOne(ctx, m)

	if err != nil {
		return nil, errors.Wrap(err, "PostMongoRepo.Create.InsertOne")
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	result := &models.Post{}

	if err = coll.FindOne(ctx, bson.D{{"_id", res.InsertedID}}).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repo) getCollection() *mongo.Collection {
	return r.dbClient.Database(dbName).Collection(collectionName)
}
