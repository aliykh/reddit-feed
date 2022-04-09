package posts

import (
	"context"
	"github.com/aliykh/reddit-feed/internal/posts/models"
)

type UseCase interface {
	Create(context.Context, *models.Post) (*models.Post, error)
}
