package repository

import (
	"context"

	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/auth"
)

type Auth interface {
	Create(ctx context.Context, token auth.Token) error
	Find(ctx context.Context, tokenID string) (*auth.Token, error)
	Update(ctx context.Context, token auth.Token) error
	Delete(ctx context.Context, tokenID string) error
}
