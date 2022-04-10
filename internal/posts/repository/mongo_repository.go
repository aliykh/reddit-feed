package repository

import (
	"context"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

const (
	dbName         = "reddit-feed"
	collectionName = "posts"
)

type Repository interface {
	Create(context.Context, *models.Post) (*models.Post, error)
	GenerateFeeds(context.Context, *pagination.Query) (*models.Feed, error)
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

func (r *repo) GenerateFeeds(ctx context.Context, query *pagination.Query) (*models.Feed, error) {

	// db.posts.aggregate([{$match: {promoted: true}},{$sample: {size: 2}}]) used to randomly selecting promoted posts in the size of 2

	coll := r.getCollection()

	totalCount, err := coll.CountDocuments(ctx, bson.D{{"promoted", false}})

	if err != nil {
		r.logger.Default().Error("Generate feeds count docs", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "GenerateFeeds.CountDocuments")
	}

	result := make([]*models.Post, 0, 27) // max 27

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	filter := bson.D{{"promoted", false}}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetSkip(int64(query.GetOffset())).SetLimit(int64(query.Size))

	cur, err := coll.Find(ctx, filter, opts)

	if err != nil {
		r.logger.Default().Error("generate feeds find", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "GenerateFeeds.Find")
	}

	err = cur.All(context.TODO(), &result)
	if err != nil {
		r.logger.Default().Error("generate feeds cursor.All", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "GenerateFeeds.cursor.all")
	}

	matchStage := bson.D{{"$match", bson.D{{"promoted", true}}}}
	sampleStage := bson.D{{"$sample", bson.D{{"size", 2}}}}

	cur, err = coll.Aggregate(ctx, mongo.Pipeline{matchStage, sampleStage})

	if err != nil {
		r.logger.Default().Error("generatefeeds.collection.aggregate", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "GenerateFeeds.GetPromotedPosts")
	}

	var promotedPosts []*models.Post

	err = cur.All(ctx, &promotedPosts)
	if err != nil {
		r.logger.Default().Error("generate feeds promoted cursor.All", zap.String("err", err.Error()))
		return nil, errors.Wrap(err, "GenerateFeeds.promoted.cursor.all")
	}

	if len(promotedPosts) > 0 && len(result) >= 3 && !*result[0].NSFW && !*result[1].NSFW {
		result = append(result, &models.Post{})
		copy(result[2:], result[1:])
		result[1] = promotedPosts[:1][0]
		promotedPosts = promotedPosts[1:]
	}

	if len(promotedPosts) > 0 && len(result) > 16 && !*result[14].NSFW && !*result[15].NSFW {
		result = append(result, &models.Post{})
		copy(result[15:], result[14:])
		result[15] = promotedPosts[:1][0]
	}

	return &models.Feed{
		TotalCount: totalCount,
		TotalPages: pagination.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       len(result),
		HasMore:    pagination.GetHasMore(query.GetPage(), int(totalCount), query.GetSize()),
		Posts:      result,
	}, nil
}

func (r *repo) getCollection() *mongo.Collection {
	return r.dbClient.Database(dbName).Collection(collectionName)
}
