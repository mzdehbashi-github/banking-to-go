package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GenerateEllipticKeyPair() (*ecdsa.PrivateKey, error) {
	// Generate a private key using the P-256 elliptic curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func SavePrivateKeyToFile(privateKey *ecdsa.PrivateKey, filePath string) error {
	// Serialize the private key
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	// Write the private key to a file
	err = os.WriteFile(filePath, keyBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func SavePublicKeyToFile(publicKey *ecdsa.PublicKey, filePath string) error {
	// Serialize the public key
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	// Write the public key to a file
	err = os.WriteFile(filePath, keyBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../..")

	// Generate a new elliptic curve key pair
	privateKey, err := GenerateEllipticKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		os.Exit(1)
	}

	// Save the private key to a file
	err = SavePrivateKeyToFile(privateKey, filepath.Join(root, "private.pem"))
	if err != nil {
		fmt.Println("Error saving private key:", err)
		os.Exit(1)
	}
	fmt.Println("Private key saved to private.pem")

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Save the public key to a file
	err = SavePublicKeyToFile(publicKey, filepath.Join(root, "public.pem"))
	if err != nil {
		fmt.Println("Error saving public key:", err)
		os.Exit(1)
	}
	fmt.Println("Public key saved to public.pem")
}
