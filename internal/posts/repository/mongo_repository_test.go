package repository

import (
	"context"
	"errors"
	logr "github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/driver/db/mock"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

// posts - for testing purposes
var posts = []*models.Post{
	{
		Title:     "title1",
		Author:    "t2_author1",
		Link:      "www.example.com",
		Subreddit: "/r/nsfw",
		Score:     new(int),
		Promoted:  new(bool),
		NSFW:      new(bool),
	},
	{
		Title:     "title2",
		Author:    "t2_author2",
		Link:      "www.example.com",
		Subreddit: "/r/nsfw2",
		Score:     new(int),
		Promoted:  new(bool),
		NSFW:      new(bool),
	},
}

func TestRepo_Create(t *testing.T) {
	var logger = logr.NewFactory(logr.Mock, "test")

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coll := mock.NewMockCollection(ctrl)

	repo := New(logger, coll)

	t.Run("ok", func(t *testing.T) {

		p := &models.Post{
			Title:     "title",
			Author:    "t2_author",
			Link:      "www.example.com",
			Subreddit: "/r/nsfw",
			Score:     new(int),
			Promoted:  new(bool),
			NSFW:      new(bool),
		}

		objcId := primitive.NewObjectID()
		coll.EXPECT().InsertOne(gomock.Any(), gomock.Eq(p)).Return(&mongo.InsertOneResult{InsertedID: objcId}, nil)
		coll.EXPECT().FindOne(gomock.Any(), bson.D{{"_id", objcId}}, gomock.Eq(&models.Post{})).Return(nil).SetArg(2, *p)

		result, err := repo.Create(context.Background(), p)

		require.NoError(t, err)
		require.Equal(t, p, result)

	})

	t.Run("fails", func(t *testing.T) {

		coll.EXPECT().InsertOne(gomock.Any(), gomock.Nil()).Return(nil, mongo.ErrNilDocument)

		_, err := repo.Create(context.Background(), nil)
		require.True(t, errors.Is(err, mongo.ErrNilDocument))
	})

	t.Run("fails-findone", func(t *testing.T) {

		p := &models.Post{}

		objcId := primitive.NewObjectID()
		coll.EXPECT().InsertOne(gomock.Any(), gomock.Eq(p)).Return(&mongo.InsertOneResult{InsertedID: objcId}, nil)

		coll.EXPECT().FindOne(gomock.Any(), bson.D{{"_id", objcId}}, gomock.Eq(&models.Post{})).Return(mongo.ErrNoDocuments)

		_, err := repo.Create(context.Background(), p)
		require.True(t, errors.Is(err, mongo.ErrNoDocuments))

	})

}

func TestRepo_CountDocuments(t *testing.T) {

	var logger = logr.NewFactory(logr.Mock, "test")

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coll := mock.NewMockCollection(ctrl)

	repo := New(logger, coll)

	t.Run("ok", func(t *testing.T) {

		filter := bson.D{{"promoted", false}}

		coll.EXPECT().CountDocuments(gomock.Any(), gomock.Eq(filter)).Return(int64(20), nil)

		c, err := repo.CountDocuments(context.Background(), filter)

		require.NoError(t, err)
		require.Equal(t, c, int64(20))

	})

	t.Run("error", func(t *testing.T) {

		filter := bson.D{{"promoted", false}}

		coll.EXPECT().CountDocuments(gomock.Any(), gomock.Eq(filter)).Return(int64(0), mongo.ErrNilDocument)

		c, err := repo.CountDocuments(context.Background(), filter)

		require.Equal(t, c, int64(0))
		require.True(t, errors.Is(err, mongo.ErrNilDocument))
	})

}

func TestRepo_Aggregate(t *testing.T) {

	var logger = logr.NewFactory(logr.Mock, "test")

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coll := mock.NewMockCollection(ctrl)

	repo := New(logger, coll)

	t.Run("ok", func(t *testing.T) {

		matchStage := bson.D{{"$match", bson.D{{"promoted", true}}}}
		sampleStage := bson.D{{"$sample", bson.D{{"size", 2}}}}

		coll.EXPECT().Aggregate(gomock.Any(), gomock.Eq(mongo.Pipeline([]bson.D{matchStage, sampleStage})), gomock.Any()).Return(nil).SetArg(2, posts)

		result, err := repo.Aggregate(context.Background(), matchStage, sampleStage)

		require.NoError(t, err)
		require.Equal(t, result, posts)
	})

	t.Run("error", func(t *testing.T) {

		matchStage := bson.D{{"$match", bson.D{{"promoted", true}}}}
		sampleStage := bson.D{{"$sample", bson.D{{"size", 2}}}}

		coll.EXPECT().Aggregate(gomock.Any(), gomock.Eq(mongo.Pipeline([]bson.D{matchStage, sampleStage})), gomock.Any()).Return(mongo.CommandError{})
		result, err := repo.Aggregate(context.Background(), matchStage, sampleStage)

		require.Empty(t, result)
		require.True(t, errors.As(err, &mongo.CommandError{}))

	})

}

func TestRepo_FindAll(t *testing.T) {

	var logger = logr.NewFactory(logr.Mock, "test")

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coll := mock.NewMockCollection(ctrl)

	repo := New(logger, coll)

	t.Run("ok", func(t *testing.T) {

		query := &pagination.Query{
			Size: 25,
			Page: 0,
		}

		filter := bson.D{{"promoted", false}}
		opts := options.Find().SetSort(bson.D{{"score", -1}}).SetSkip(int64(query.GetOffset())).SetLimit(int64(query.GetSize()))

		coll.EXPECT().Find(gomock.Any(), gomock.Eq(filter), gomock.Eq(&[]*models.Post{}), opts).Return(nil).SetArg(2, posts)

		result, err := repo.FindAll(context.Background(), filter, query)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, posts, result)

	})

	t.Run("error", func(t *testing.T) {

		coll.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mongo.ErrNoDocuments)
		result, err := repo.FindAll(context.Background(), bson.D{}, &pagination.Query{})
		require.Empty(t, result)
		require.True(t, errors.Is(err, mongo.ErrNoDocuments))
	})

}
