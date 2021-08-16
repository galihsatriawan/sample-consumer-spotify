package main

import (
	spotifyHandler "github.com/galihsatriawan/sample-consumer-spotify/internal/app/spotify/delivery/http"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/repository"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/spotify"
	redisStorage "github.com/galihsatriawan/sample-consumer-spotify/internal/storage/redis"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

var rds *redis.Client
var authStorage repository.Auth

func init() {
	/** Initiate Connection Redis **/
	rds = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

}
func main() {
	e := echo.New()
	authStorage = redisStorage.NewAuth(rds)
	spotifyUc := spotify.NewSpotifyUsecase(authStorage)
	spotifyHandler.NewSpotifyHandler(e, spotifyUc)
	e.Start(":4000")
}
