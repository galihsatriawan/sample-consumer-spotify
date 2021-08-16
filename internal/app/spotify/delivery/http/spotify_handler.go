package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/galihsatriawan/sample-consumer-spotify/config"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/auth"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/spotify"
	"github.com/galihsatriawan/sample-consumer-spotify/middleware"
	"github.com/labstack/echo"
)

var (
	SESSION_LOGIN = spotify.LoginResponse{}
)

type SpotifyHandler struct {
	e              *echo.Echo
	spotifyUsecase spotify.SpotifyUsecaseProto
}

func NewSpotifyHandler(e *echo.Echo, spotifyUsecase spotify.SpotifyUsecaseProto, authMiddleware middleware.AuthMiddleware) SpotifyHandler {
	spotifyHandler := &SpotifyHandler{e: e, spotifyUsecase: spotifyUsecase}

	groupsSpotify := e.Group("spotify")

	groupsSpotify.GET("/callback", spotifyHandler.CallbackHandler)
	groupsSpotify.GET("/auth", spotifyHandler.AuthHandler)
	groupsSpotify.GET("/login", spotifyHandler.LoginHandler)
	groupsSpotify.GET("/me", spotifyHandler.ProfileHandler, authMiddleware.Auth)
	groupsSpotify.GET("/playlists", spotifyHandler.PlaylistsHandler, authMiddleware.Auth)
	groupsSpotify.GET("/new_releases", spotifyHandler.NewReleaseHandler, authMiddleware.Auth)
	return *spotifyHandler
}

func (h *SpotifyHandler) CallbackHandler(e echo.Context) error {
	config.CODE = e.QueryParams().Get("code")
	return e.JSON(http.StatusOK, spotify.Response{
		Status:  true,
		Message: "Success",
		Data:    config.CODE,
	})
}
func (h *SpotifyHandler) AuthHandler(e echo.Context) error {
	params := url.Values{
		"client_id":     {config.CLIENT_ID},
		"client_secret": {config.CLIENT_SECRET},
		"scope":         {config.SCOPE},
		"response_type": {"code"},
		"redirect_uri":  {config.REDIRECT_URL},
		"state":         {config.STATE},
	}
	reqAuthorizeUrl := fmt.Sprintf("%v?%v", config.AUTHORIZE_URL, params.Encode())

	return e.Redirect(http.StatusSeeOther, reqAuthorizeUrl)
}
func (h *SpotifyHandler) ProfileHandler(e echo.Context) error {
	req, err := http.NewRequest(http.MethodGet, config.PROFILE_URL, nil)
	if err != nil {
		panic(err)
	}
	authorization := fmt.Sprintf("Bearer %v", SESSION_LOGIN.AccessToken)
	req.Header.Add("Authorization", authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var myProfile map[string]interface{}
	err = json.Unmarshal(bodyResp, &myProfile)
	if err != nil {
		panic(err)
	}
	return e.JSON(http.StatusOK, myProfile)
}

func (h *SpotifyHandler) LoginHandler(e echo.Context) error {

	bodyRequest := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {config.CODE},
		"redirect_uri": {config.REDIRECT_URL},
	}

	req, err := http.NewRequest(http.MethodPost, config.LOGIN_URL, strings.NewReader(bodyRequest.Encode()))

	authorization := fmt.Sprintf("%v:%v", config.CLIENT_ID, config.CLIENT_SECRET)
	encodedClient := base64.StdEncoding.EncodeToString([]byte(authorization))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", encodedClient))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var bodyObject spotify.LoginResponse
	err = json.Unmarshal(body, &bodyObject)
	if err != nil {
		panic(err)
	}
	SESSION_LOGIN = bodyObject

	// Save to redis
	credential := spotify.FormatFromSpotifyResponse(bodyObject)
	h.spotifyUsecase.SaveToken(credential)

	return e.JSON(http.StatusOK, credential)
}
func (h *SpotifyHandler) PlaylistsHandler(e echo.Context) error {
	req, err := http.NewRequest(http.MethodGet, config.PLAYLIST_URL, nil)
	if err != nil {
		panic(err)
	}

	credential := e.Get("token").(*auth.Token)

	authorization := fmt.Sprintf("Bearer %v", credential.AccessToken)
	req.Header.Add("Authorization", authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var playlistObject map[string]interface{}
	err = json.Unmarshal(bodyResp, &playlistObject)
	if err != nil {
		panic(err)
	}
	return e.JSON(http.StatusOK, playlistObject)
}

func (h *SpotifyHandler) NewReleaseHandler(e echo.Context) error {
	req, err := http.NewRequest(http.MethodGet, config.NEW_RELEASE_URL, nil)
	if err != nil {
		panic(err)
	}
	authorization := fmt.Sprintf("Bearer %v", SESSION_LOGIN.AccessToken)
	req.Header.Add("Authorization", authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var tracksObject map[string]interface{}
	err = json.Unmarshal(bodyResp, &tracksObject)
	if err != nil {
		panic(err)
	}
	return e.JSON(http.StatusOK, tracksObject)
}
