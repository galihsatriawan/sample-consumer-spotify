package auth

import (
	"encoding/json"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	Spotify      Spotify
}

func (t *Token) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Token) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}
