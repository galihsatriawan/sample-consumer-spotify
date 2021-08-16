package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var (
	CLIENT_ID     = ""
	CLIENT_SECRET = ""
	AUTHORIZE_URL = ""
	SCOPE         = "user-read-private user-read-email"
	REDIRECT_URL  = ""
	CODE          = ""
	LOGIN_URL     = ""
	STATE         = ""
	PROFILE_URL   = ""
	SESSION_LOGIN = LoginResponse{}
)

func init() {
	viper.SetConfigName("dev")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	CLIENT_ID = viper.GetString("client_id")
	CLIENT_SECRET = viper.GetString("client_secret")
	AUTHORIZE_URL = viper.GetString("authorize_url")
	REDIRECT_URL = viper.GetString("redirect_uri")
	LOGIN_URL = viper.GetString("token_url")
	STATE = viper.GetString("state")
	PROFILE_URL = viper.GetString("profile_url")
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func CallbackHandler(e echo.Context) error {
	CODE = e.QueryParams().Get("code")
	return e.JSON(http.StatusOK, Response{
		Status:  true,
		Message: "Success",
		Data:    CODE,
	})
}
func AuthHandler(e echo.Context) error {
	params := url.Values{
		"client_id":     {CLIENT_ID},
		"client_secret": {CLIENT_SECRET},
		"scope":         {SCOPE},
		"response_type": {"code"},
		"redirect_uri":  {REDIRECT_URL},
		"state":         {STATE},
	}
	reqAuthorizeUrl := fmt.Sprintf("%v?%v", AUTHORIZE_URL, params.Encode())

	return e.Redirect(http.StatusSeeOther, reqAuthorizeUrl)
}
func ProfileHandler(e echo.Context) error {
	req, err := http.NewRequest(http.MethodGet, PROFILE_URL, nil)
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

type LoginRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func LoginHandler(e echo.Context) error {
	// logRequest := LoginRequest{
	// 	GrantType:   "authorization_code",
	// 	Code:        CODE,
	// 	RedirectUri: REDIRECT_URL,
	// }
	bodyRequest := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {CODE},
		"redirect_uri": {REDIRECT_URL},
	}
	// jsonBody, err := json.Marshal(logRequest)
	// if err != nil {
	// 	panic(err)
	// }

	req, err := http.NewRequest(http.MethodPost, LOGIN_URL, strings.NewReader(bodyRequest.Encode()))

	authorization := fmt.Sprintf("%v:%v", CLIENT_ID, CLIENT_SECRET)
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
	var bodyObject LoginResponse
	err = json.Unmarshal(body, &bodyObject)
	if err != nil {
		panic(err)
	}
	SESSION_LOGIN = bodyObject
	return e.JSON(http.StatusOK, SESSION_LOGIN)
}
func main() {
	e := echo.New()
	e.GET("/callback", CallbackHandler)
	e.GET("/auth", AuthHandler)
	e.GET("/login", LoginHandler)
	e.GET("/me", ProfileHandler)
	e.Start(":4000")
}
