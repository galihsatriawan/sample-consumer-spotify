package auth

// TokenID is refresh token/refresh_token
func CacheHashKey(tokenID string) string {
	return "app:auth:" + tokenID
}
func CacheHashField() string {
	return "data"
}
