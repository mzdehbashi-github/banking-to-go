package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateEllipticKeyPair() (*ecdsa.PrivateKey, error) {
	// Generate a private key using the P-256 elliptic curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func PrivateKeyToString(privateKey *ecdsa.PrivateKey) (string, error) {
	// Serialize the private key
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	// Encode the private key to PEM format
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	privateKeyPEM := pem.EncodeToMemory(block)

	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)
	return privateKeyBase64, nil
}

func PublicKeyToString(publicKey *ecdsa.PublicKey) (string, error) {
	// Serialize the public key
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	// Encode the public key to PEM format
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	publicKeyPEM := pem.EncodeToMemory(block)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)

	return publicKeyBase64, nil
}

func main() {
	// Generate a new elliptic curve key pair
	privateKey, err := GenerateEllipticKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		os.Exit(1)
	}

	// Convert private key to string
	privateKeyStr, err := PrivateKeyToString(privateKey)
	if err != nil {
		fmt.Println("Error converting private key to string:", err)
		os.Exit(1)
	}

	// Print private key
	fmt.Println("Private key:")
	fmt.Println(privateKeyStr)

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Convert public key to string
	publicKeyStr, err := PublicKeyToString(publicKey)
	if err != nil {
		fmt.Println("Error converting public key to string:", err)
		os.Exit(1)
	}

	// Print public key
	fmt.Println("Public key:")
	fmt.Println(publicKeyStr)
}
