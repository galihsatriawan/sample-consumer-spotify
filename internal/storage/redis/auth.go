package redis

import (
	"context"
	"fmt"
	"time"

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
func (a Auth) Create(ctx context.Context, token auth.Token) error {
	// Access Token
	if _, err := a.rds.HSetNX(ctx, auth.CacheHashKey(token.AccessToken), auth.CacheHashField(), &token).Result(); err != nil {
		return fmt.Errorf("create: redis error: %w", err)
	}
	a.rds.Expire(ctx, auth.CacheHashKey(token.AccessToken), time.Duration(auth.EXPIRED_IN*int(time.Second)))
	// Refresh Token
	if _, err := a.rds.HSetNX(ctx, auth.CacheHashKey(token.RefreshToken), auth.CacheHashField(), &token).Result(); err != nil {
		return fmt.Errorf("create: redis error: %w", err)
	}
	a.rds.Expire(ctx, auth.CacheHashKey(token.AccessToken), time.Duration(auth.EXPIRED_IN*int(time.Second)))
	return nil
}
func (a Auth) Find(ctx context.Context, tokenID string) (*auth.Token, error) {
	result, err := a.rds.HGet(ctx, auth.CacheHashKey(tokenID), auth.CacheHashField()).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("find: redis error: %w", err)
	}
	if result == "" {
		return nil, fmt.Errorf("find: not found")
	}

	token := &auth.Token{}
	if err := token.UnmarshalBinary([]byte(result)); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}

	return token, nil
}
func (a Auth) Update(ctx context.Context, token auth.Token) error {
	panic("do something")
}
func (a Auth) Delete(ctx context.Context, tokenID string) error {
	panic("do something")
}
