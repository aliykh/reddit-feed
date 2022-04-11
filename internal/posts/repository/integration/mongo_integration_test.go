package integration

import (
	"context"
	logr "github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/config"
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
)

func TestPosts_Create(t *testing.T) {

	// load application configurations
	cfg, err := config.Load("../../../../config/local.yml", nil)
	if err != nil {
		log.Fatalf("failed to load application configuration: %s\n", err)
	}

	clientOpts := options.Client()
	clientOpts.ApplyURI(cfg.MongoAddr)
	clientOpts.SetMaxPoolSize(5)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err)

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

	repo := repository.New(logger, client, DatabaseName)

	createdPost, err := repo.Create(context.Background(), m)

	require.NoError(t, err)
	require.NotEmpty(t, createdPost.Id)
	require.Equal(t, createdPost.Score, m.Score)
}