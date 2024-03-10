package token

import (
	"errors"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker() (TokenMaker, error) {
	secretKey := os.Getenv("PASETO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "YELLOW SUBMARINE, BLACK WIZARDRY"
		// return nil, errors.New("paseto secret key is empty")
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(secretKey),
	}, nil
}

func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	claims := paseto.JSONToken{
		// Audience:   "user123",
		Subject:    username,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(duration), // Token expires in 24 hours
	}

	// Create a PASETO token
	token, err := pm.paseto.Encrypt(pm.symmetricKey, claims, nil)
	if err != nil {
		return "", errors.New("error creating token")
	}
	return token, nil

}

func (pm *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	var token paseto.JSONToken
	err := pm.paseto.Decrypt(tokenString, pm.symmetricKey, &token, nil)
	if err != nil {
		return nil, errors.New("error in verifying token")
	}
	if token.Expiration.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}
	payload := &Payload{
		Username: token.Subject,
	}
	return payload, nil

}
