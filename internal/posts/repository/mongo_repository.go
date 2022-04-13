package repository

import (
	"context"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/driver/db"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

const (
	collectionName = "posts"
)

type Repository interface {
	Create(context.Context, *models.Post) (*models.Post, error)
	CountDocuments(ctx context.Context, filter bson.D) (int64, error)
	FindAll(ctx context.Context, filter bson.D, query *pagination.Query) ([]*models.Post, error)
	Aggregate(ctx context.Context, stages ...bson.D) ([]*models.Post, error)
}

type repo struct {
	logger     *log.Factory
	collection db.Collection
}

func New(logger *log.Factory, collection db.Collection) *repo {
	return &repo{
		logger:     logger,
		collection: collection,
	}
}

func (r *repo) Create(ctx context.Context, m *models.Post) (*models.Post, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, m)

	if err != nil {
		return nil, errors.Wrap(err, "PostMongoRepo.Create.InsertOne")
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	result := &models.Post{}

	if err = r.collection.FindOne(ctx, bson.D{{"_id", res.InsertedID}}, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repo) CountDocuments(ctx context.Context, filter bson.D) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	totalCount, err := r.collection.CountDocuments(ctx, filter)

	if err != nil {
		r.logger.Default().Error("Generate feeds count docs", zap.String("err", err.Error()))
		return 0, errors.Wrap(err, "GenerateFeeds.CountDocuments")
	}

	return totalCount, nil
}

func (r *repo) FindAll(ctx context.Context, filter bson.D, query *pagination.Query) ([]*models.Post, error) {

	result := make([]*models.Post, 0, query.GetSize()+2) // max 27

	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetSkip(int64(query.GetOffset())).SetLimit(int64(query.GetSize()))

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	err := r.collection.Find(ctx, filter, &result, opts)

	if err != nil {
		r.logger.Default().Error("FindAll.Find", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "FindAll.Find")
	}

	return result, nil
}

func (r *repo) Aggregate(ctx context.Context, stages ...bson.D) ([]*models.Post, error) {

	var result []*models.Post

	err := r.collection.Aggregate(ctx, stages, &result)

	if err != nil {
		r.logger.Default().Error("Aggreegate.Collections.GetPosts", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "Aggreegate.GetPosts")
	}

	return result, nil
}
