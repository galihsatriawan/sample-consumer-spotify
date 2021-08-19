package spotify

import (
	"context"

	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/auth"
	"github.com/galihsatriawan/sample-consumer-spotify/internal/domain/repository"
)

type spotifyUsecase struct {
	authRepository repository.Auth
}
type SpotifyUsecaseProto interface {
	SaveToken(credential CredentialEntity)
}

func NewSpotifyUsecase(authRepository repository.Auth) SpotifyUsecaseProto {
	return spotifyUsecase{authRepository: authRepository}
}
func (uc spotifyUsecase) SaveToken(credential CredentialEntity) {
	getToken := auth.Token{
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		Spotify: auth.Spotify{
			Scope:     credential.Scope,
			TokenType: credential.TokenType,
		},
	}
	uc.authRepository.Create(context.Background(), getToken)
}
