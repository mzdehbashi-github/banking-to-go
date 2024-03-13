package util

import (
	"os"
)

// Config stores all configuration of the application.
type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
	PublicKey     string
	PrivateKey    string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() Config {
	var C Config
	C.DBDriver = os.Getenv("DB_DRIVER")
	C.DBSource = os.Getenv("DB_SOURCE")
	C.ServerAddress = os.Getenv("SERVER_ADDRESS")
	C.PublicKey = os.Getenv("PUBLIC_KEY")
	C.PrivateKey = os.Getenv("PRIVATE_KEY")
	return C
}
