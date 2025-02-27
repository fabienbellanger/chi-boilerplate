package pkg

import "github.com/spf13/viper"

// Config represents the configuration of the application in .env file
type Config struct {
	// Application environment (development, production or test)
	AppEnv string

	// Application name
	AppName string
}

func NewConfig() (*Config, error) {
	// Read .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:  viper.GetString("APP_ENV"),
		AppName: viper.GetString("APP_NAME"),
	}, nil
}
