package auth

type Token struct {
	ID   string // Refresh token
	User User
}
