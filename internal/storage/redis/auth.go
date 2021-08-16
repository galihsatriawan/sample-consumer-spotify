package storage

import (
	"context"

	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/auth"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/repository"
	"github.com/go-redis/redis/v8"
)

type Auth struct {
	rds *redis.Client
}

func NewAuth(rds *redis.Client) repository.Auth {
	return &Auth{rds: rds}
}
func (a Auth) Create(ctx context.Context, token auth.Token) {
	panic("do something")
}
func (a Auth) Find(ctx context.Context, tokenID string) {
	panic("do something")
}
func (a Auth) Update(ctx context.Context, token auth.Token) {
	panic("do something")
}
func (a Auth) Delete(ctx context.Context, tokenID string) {
	panic("do something")
}
