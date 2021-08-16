package repository

import (
	"context"

	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/auth"
)

type Auth interface {
	Create(ctx context.Context, token auth.Token)
	Find(ctx context.Context, tokenID string)
	Update(ctx context.Context, token auth.Token)
	Delete(ctx context.Context, tokenID string)
}
