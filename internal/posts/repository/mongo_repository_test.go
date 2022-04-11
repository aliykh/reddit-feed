package repository

import (
	"context"
	"encoding/json"
	"fmt"
	logr "github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"github.com/stretchr/testify/require"
	"github.com/strikesecurity/strikememongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	// collection constants
	postsCollectionName = "posts"
	databaseName        = "reddit-feed-test"
)

var (

	// collections variables
	postsCollection *mongo.Collection
	mongoURI        = ""
	database        *mongo.Database
	mongoClient     *mongo.Client
)

func TestMain(m *testing.M) {
	mongoServer, err := strikememongo.StartWithOptions(&strikememongo.Options{MongoVersion: "4.2.0", ShouldUseReplica: true})
	if err != nil {
		log.Fatal(err)
	}
	mongoURI = mongoServer.URIWithRandomDB()

	defer mongoServer.Stop()

	setup()
	m.Run()
}

func setup() {
	startApplication()
	createCollections()
	cleanup()
}

// startApplication initializes the engine and the necessary components for the (test) service to work
func startApplication() {
	// Initialize Database (memongodb)
	dbClient, ctx, err := initDB()
	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	mongoClient = dbClient
	database = dbClient.Database(databaseName)
}

func initDB() (client *mongo.Client, ctx context.Context, err error) {
	uri := fmt.Sprintf("%s%s", mongoURI, "?retryWrites=false")
	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return
	}

	return
}

// createCollections cretaes the necessary collections to be used during tests
func createCollections() {
	err := database.CreateCollection(context.Background(), postsCollectionName)
	if err != nil {
		fmt.Printf("error creating collection: %s", err.Error())
	}

	postsCollection = database.Collection(postsCollectionName)
}

// ----------------------------
// 		TEAR DOWN FUNCTION
// ----------------------------
func cleanup() {
	postsCollection.DeleteMany(context.Background(), bson.M{})
}

// ----------------------------
// 		TEST THE ACTUAL REPOSITORY
// ----------------------------
func TestRepo_Create(t *testing.T) {

	m := &models.Post{
		Title:     "title",
		Subreddit: "/r/subreddit",
		Content:   "content text",
		Promoted:  new(bool),
		NSFW:      new(bool),
		Score:     new(int),
	}
	*m.Score = 100
	m.GenerateAuthorName()

	logger := logr.NewFactory(logr.Mock, "test")

	repo := New(logger, mongoClient, databaseName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createdPost, err := repo.Create(ctx, m)

	require.NoError(t, err)
	require.NotEmpty(t, createdPost.Id)
	require.Equal(t, *createdPost.Score, 100)

}

func TestRepo_GenerateFeeds(t *testing.T) {
	cleanup()

	posts := prepareDataInDbForGeneratePosts(t, "sample.json")

	// ----------------------------
	// 		TEST THE REQUIREMENTS POINTS 3-4
	// ----------------------------
	logger := logr.NewFactory(logr.Mock, "test")

	repo := New(logger, mongoClient, databaseName)

	postsInDb, err := repo.GenerateFeeds(context.Background(), &pagination.Query{
		Size: 25,
		Page: 0,
	})

	require.NoError(t, err)
	require.Equal(t, len(posts), len(postsInDb.Posts))
	require.Equal(t, true, *postsInDb.Posts[1].Promoted)
	require.Equal(t, true, *postsInDb.Posts[15].Promoted)

}

func TestRepo_GenerateFeeds_NSFW(t *testing.T) {
	cleanup()

	posts := prepareDataInDbForGeneratePosts(t, "sample-2.json")

	// ----------------------------
	// 		TEST THE REQUIREMENTS POINTS 5
	// ----------------------------
	logger := logr.NewFactory(logr.Mock, "test")

	repo := New(logger, mongoClient, databaseName)

	postsInDb, err := repo.GenerateFeeds(context.Background(), &pagination.Query{
		Size: 25,
		Page: 0,
	})

	require.NoError(t, err)
	require.Equal(t, len(posts)-1, len(postsInDb.Posts)) // 21 - minus one of the promoted posts due to the fact that one of its adjacent posts is NSFW type
	require.Equal(t, false, *postsInDb.Posts[1].Promoted)
	require.Equal(t, true, *postsInDb.Posts[15].Promoted)

}

func prepareDataInDbForGeneratePosts(t *testing.T, fileName string) []models.Post {
	// ----------------------------
	// 		DATA PREP IN THE DB
	// ----------------------------
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)

	require.NoError(t, err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	var posts = make([]models.Post, 0, 30)

	err = json.Unmarshal(data, &posts)
	require.NoError(t, err)

	bulkWriteData := make([]mongo.WriteModel, 0, 30)
	for _, v := range posts {
		op := mongo.NewInsertOneModel()
		op.SetDocument(v)
		bulkWriteData = append(bulkWriteData, op)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	coll := mongoClient.Database(databaseName).Collection(collectionName)
	result, err := coll.BulkWrite(ctx, bulkWriteData)
	require.NoError(t, err)
	require.Equal(t, int64(len(posts)), result.InsertedCount)

	return posts
}

func TestGenerateAuthorName(t *testing.T) {

	m := &models.Post{}
	m.GenerateAuthorName()

	require.NotEmpty(t, m.Author)
	require.True(t, strings.HasPrefix(m.Author, "t2_"))

}
