package usecase

import (
	"context"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/internal/posts/repository"
	"github.com/aliykh/reddit-feed/pkg/pagination"
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

func (p *postsUC) GenerateFeeds(ctx context.Context, page *pagination.Query) (*models.Feed, error)  {
	return p.repo.GenerateFeeds(ctx, page)
}