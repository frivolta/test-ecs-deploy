package util

import (
	"firebase.google.com/go/auth"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Auth0Domain   string `mapstructure:"AUTH0_DOMAIN"`
	Auth0Iss      string `mapstructure:"AUTH0_ISS"`
	Auth0Audience string `mapstructure:"AUTH0_AUDIENCE"`
	TestBearer    string `mapstructure:"TEST_BEARER"`
	TestUUID      string `mapstructure:"TEST_UUID"`
	FirebaseWeb   string `mapstructure:"FIREBASE_WEB"`
	GinMode       string `mapstructure:"GIN_MODE"`
	AuthClient    *auth.Client
}

func LoadConfig(path string, serviceKeyPath string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	// Setup firebase
	config.AuthClient = SetupFirebase(serviceKeyPath)

	return
}
