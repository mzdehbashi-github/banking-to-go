package token

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func loadPublicKeyFromString(keyString string) (*ecdsa.PublicKey, error) {

	key, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPublicKey, ok := pubKey.(*ecdsa.PublicKey)
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

func loadPrivateKeyFromString(keyString string) (*ecdsa.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func NewJWTMaker(privateKeyString, publicKeyString string) (TokenMaker, error) {
	privateKey, err := loadPrivateKeyFromString(privateKeyString)
	if err != nil {
		log.Fatal("error in parsing private key: ", err)
	}

	publicKey, err := loadPublicKeyFromString(publicKeyString)
	if err != nil {
		log.Fatal("error in parsing public key: ", err)
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
