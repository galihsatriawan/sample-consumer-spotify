package auth

var (
	EXPIRED_IN            = 3600
	EXPIRED_REFRESH_TOKEN = 7200
)

type Spotify struct {
	Scope     string
	TokenType string
}
