package usecase

import (
	"context"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/internal/posts/repository"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type postsUC struct {
	logger *log.Factory
	repo   repository.Repository
}

func New(logger *log.Factory, repo repository.Repository) *postsUC {
	return &postsUC{
		logger: logger,
		repo:   repo,
	}
}

func (p *postsUC) Create(ctx context.Context, model *models.Post) (*models.Post, error) {
	model.GenerateAuthorName()
	return p.repo.Create(ctx, model)
}

func (p *postsUC) GenerateFeeds(ctx context.Context, query *pagination.Query) (*models.Feed, error) {

	totalCount, err := p.repo.CountDocuments(ctx, bson.D{{"promoted", false}})

	if err != nil {
		return nil, err
	}

	posts, err := p.repo.FindAll(ctx, bson.D{{"promoted", false}}, query)
	if err != nil {
		return nil, err
	}

	matchStage := bson.D{{"$match", bson.D{{"promoted", true}}}}
	sampleStage := bson.D{{"$sample", bson.D{{"size", 2}}}}

	promotedPosts, err := p.repo.Aggregate(ctx, matchStage, sampleStage)

	if err != nil {
		return nil, err
	}

	if len(promotedPosts) > 0 && len(posts) >= 3 && !*posts[0].NSFW && !*posts[1].NSFW {
		posts = append(posts, &models.Post{})
		copy(posts[2:], posts[1:])
		posts[1] = promotedPosts[:1][0]
		promotedPosts = promotedPosts[1:]
	}

	if len(promotedPosts) > 0 && len(posts) > 16 && !*posts[14].NSFW && !*posts[15].NSFW {
		posts = append(posts, &models.Post{})
		copy(posts[15:], posts[14:])
		posts[15] = promotedPosts[:1][0]
	}

	return &models.Feed{
		TotalCount: totalCount,
		TotalPages: pagination.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       len(posts),
		HasMore:    pagination.GetHasMore(query.GetPage(), int(totalCount), query.GetSize()),
		Posts:      posts,
	}, nil
}
