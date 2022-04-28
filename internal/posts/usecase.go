//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package posts

import (
	"context"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
)

type UseCase interface {
	Create(context.Context, *models.Post) (*models.Post, error)
	GenerateFeeds(context.Context, *pagination.Query) (*models.Feed, error)
}
