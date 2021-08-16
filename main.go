package main

import (
	"encoding/base64"
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
	}
	reqAuthorizeUrl := fmt.Sprintf("%v?%v", AUTHORIZE_URL, params.Encode())

	return e.Redirect(http.StatusSeeOther, reqAuthorizeUrl)
}

type LoginRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
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
	fmt.Println(resp.Status)
	return e.HTML(http.StatusOK, string(body))
}
func main() {
	e := echo.New()
	e.GET("/callback", CallbackHandler)
	e.GET("/auth", AuthHandler)
	e.GET("/login", LoginHandler)
	e.Start(":4000")
}
