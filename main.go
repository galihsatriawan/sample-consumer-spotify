package main

import (
	spotifyHandler "github.com/galihsatriawan/sample-consumer-spotify/internal/app/spotify/delivery/http"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

var rds *redis.Client

func init() {
	/** Initiate Connection Redis **/
	rds = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

}
func main() {
	e := echo.New()
	spotifyHandler.NewSpotifyHandler(e)
	e.Start(":4000")
}
