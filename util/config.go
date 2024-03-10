package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or env variables.
type Config struct {
	DBDrive       string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from a file or env variables
func LoadConfig(path string) (Config, error) {
	// viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	// viper.SetConfigType("env") // json, xml
	viper.AutomaticEnv()

	// err = viper.ReadInConfig()
	// if err != nil {
	// 	return
	// }
	dbSource := viper.Get("DB_SOURCE")
	serverAddress := viper.Get("SERVER_ADDRESS")
	dbDriver := viper.Get("DB_DRIVER")
	var C Config
	C.DBDrive = dbDriver.(string)
	C.DBSource = dbSource.(string)
	C.ServerAddress = serverAddress.(string)
	return C, nil
}
