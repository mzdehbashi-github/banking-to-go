package token

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSizeSecretKey = 3
)

type JWTMaker struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func loadPublicKeyFromFile(filePath string) (*ecdsa.PublicKey, error) {
	// Read the public key file
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the raw key bytes into an elliptic curve public key
	publicKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	// Assert that the parsed key is an ECDSA public key
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key type")
	}

	return ecdsaPublicKey, nil
}

func loadPrivateKeyFromFile(filePath string) (*ecdsa.PrivateKey, error) {
	// Read the private key file
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParseECPrivateKey([]byte(keyBytes))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func NewJWTMaker() (TokenMaker, error) {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "..")
	privateKey, err := loadPrivateKeyFromFile(filepath.Join(root, "private.pem"))
	if err != nil {
		log.Fatal("error in reading private key: ", err)
	}
	publicKey, err := loadPublicKeyFromFile(filepath.Join(root, "public.pem"))
	if err != nil {
		log.Fatal(err)
	}

	return &JWTMaker{publicKey: publicKey, privateKey: privateKey}, nil

}

func (jm *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	paylaod, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, paylaod)
	return jwtToken.SignedString(jm.privateKey)

}

func (jm *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Payload{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				return nil, errors.New("invalid method")
			}

			return jm.publicKey, nil
		},
	)

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*Payload); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}
