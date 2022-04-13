package integration

import (
	"context"
	logr "github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/config"
	"github.com/aliykh/reddit-feed/internal/driver/db"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/internal/posts/repository"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"testing"
	"time"
)

const (
	DatabaseName = "reddit-feed-test"
	CollName     = "posts"
)

var (
	dbClient *mongo.Client
)

func TestMain(m *testing.M) {

	// load application configurations
	cfg, err := config.Load("../../../../config/local.yml", nil)
	if err != nil {
		log.Fatalf("failed to load application configuration: %s\n", err)
	}

	clientOpts := options.Client()
	clientOpts.ApplyURI(cfg.MongoAddr)
	clientOpts.SetMaxPoolSize(5)

	client, err := mongo.Connect(context.Background(), clientOpts)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err.Error())
	}

	dbClient = client

	m.Run()
}

func TestPosts_Create(t *testing.T) {

	m := &models.Post{
		Title:     "title",
		Subreddit: "/r/subreddit",
		Content:   "content text",
		Promoted:  new(bool),
		NSFW:      new(bool),
		Score:     new(int),
	}
	*m.Score = 1000
	m.GenerateAuthorName()

	logger := logr.NewFactory(logr.Mock, "test")

	coll := db.New(logger, dbClient, DatabaseName, CollName)
	repo := repository.New(logger, coll)

	createdPost, err := repo.Create(context.Background(), m)

	require.NoError(t, err)
	require.NotEmpty(t, createdPost.Id)
	require.Equal(t, createdPost.Score, m.Score)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	_ = dbClient.Database(DatabaseName).Collection(CollName).Drop(ctx)

}
