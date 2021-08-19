package spotify

type CredentialEntity struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func FormatFromSpotifyResponse(src LoginResponse) CredentialEntity {
	dst := CredentialEntity{}
	dst.AccessToken = src.AccessToken
	dst.RefreshToken = src.RefreshToken
	dst.Scope = src.Scope
	dst.TokenType = src.TokenType
	dst.ExpiresIn = src.ExpiresIn
	return dst
}
