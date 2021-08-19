package config

import "github.com/spf13/viper"

var (
	CLIENT_ID       = ""
	CLIENT_SECRET   = ""
	AUTHORIZE_URL   = ""
	SCOPE           = "user-read-private user-read-email"
	REDIRECT_URL    = ""
	CODE            = ""
	LOGIN_URL       = ""
	PROFILE_URL     = ""
	PLAYLIST_URL    = ""
	TRACK_URL       = ""
	NEW_RELEASE_URL = ""
	STATE           = ""
)

func init() {
	viper.SetConfigName("dev")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	/**
		Get config data
	**/
	CLIENT_ID = viper.GetString("client_id")
	CLIENT_SECRET = viper.GetString("client_secret")
	AUTHORIZE_URL = viper.GetString("authorize_url")
	REDIRECT_URL = viper.GetString("redirect_uri")
	LOGIN_URL = viper.GetString("token_url")
	STATE = viper.GetString("state")
	PROFILE_URL = viper.GetString("profile_url")
	PLAYLIST_URL = viper.GetString("playlists_url")
	NEW_RELEASE_URL = viper.GetString("new_release_url")
}
