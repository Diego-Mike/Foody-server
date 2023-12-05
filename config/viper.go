package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type EnvVariables struct {
	PORT                   string `mapstructure:"PORT"`
	DB_DRIVER              string `mapstructure:"DB_DRIVER"`
	DB_SOURCE              string `mapstructure:"DB_SOURCE"`
	GOOGLE_CLIENT_ID       string `mapstructure:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET   string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GOOGLE_REDIRECT_URI    string `mapstructure:"GOOGLE_REDIRECT_URI"`
	GOOGLE_GRANT_TYPE      string `mapstructure:"GOOGLE_GRANT_TYPE"`
	FACEBOOK_CLIENT_ID     string `mapstructure:"FACEBOOK_CLIENT_ID"`
	FACEBOOK_CLIENT_SECRET string `mapstructure:"FACEBOOK_CLIENT_SECRET"`
	FACEBOOK_REDIRECT_URI  string `mapstructure:"FACEBOOK_REDIRECT_URI"`
	REFRESH_TOKEN_KEY      string `mapstructure:"REFRESH_TOKEN_KEY"`
	REFRESH_TOKEN_TIME     string `mapstructure:"REFRESH_TOKEN_TIME"`
	ACCESS_TOKEN_KEY       string `mapstructure:"ACCESS_TOKEN_KEY"`
	ACCESS_TOKEN_TIME      string `mapstructure:"ACCESS_TOKEN_TIME"`
	API_KEY                string `mapstructure:"API_KEY"`
	SECURE_COOKIES         string `mapstructure:"SECURE_COOKIES"`
	WEBSTIE                string `mapstructure:"WEBSTIE"`
	DOMAIN                 string `mapstructure:"DOMAIN"`
}

func LoadEnv(path string) (config EnvVariables, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error unmarshaling env file", err)
	}

	// validate config
	if config.PORT == "" {
		err = errors.New("env PORT must be provided to star the server")
		return
	}

	if config.DB_DRIVER == "" {
		err = errors.New("env DB_DRIVER must be provided for db config")
		return
	}

	if config.DB_SOURCE == "" {
		err = errors.New("env DB_SOURCE must be provided to star foody_db")
		return
	}

	if config.GOOGLE_CLIENT_ID == "" {
		err = errors.New("env GOOGLE_CLIENT_ID must be provided to complete oauth flow")
		return
	}

	if config.GOOGLE_CLIENT_SECRET == "" {
		err = errors.New("env GOOGLE_CLIENT_SECRET must be provided to complete oauth flow")
		return
	}

	if config.GOOGLE_REDIRECT_URI == "" {
		err = errors.New("env GOOGLE_REDIRECT_URI must be provided to complete oauth flow")
		return
	}

	if config.GOOGLE_GRANT_TYPE == "" {
		err = errors.New("env GOOGLE_GRANT_TYPE must be provided to complete oauth flow")
		return
	}

	if config.FACEBOOK_CLIENT_ID == "" {
		err = errors.New("env FACEBOOK_CLIENT_ID must be provided to complete facebook oauth flow")
		return
	}

	if config.FACEBOOK_CLIENT_SECRET == "" {
		err = errors.New("env FACEBOOK_CLIENT_SECRET must be provided to complete facebook oauth flow")
		return
	}

	if config.FACEBOOK_REDIRECT_URI == "" {
		err = errors.New("env FACEBOOK_REDIRECT_URI must be provided to complete facebook oauth flow")
		return
	}

	if config.ACCESS_TOKEN_KEY == "" {
		err = errors.New("env ACCESS_TOKEN_KEY must be provided to complete auth flow")
		return
	}

	if config.REFRESH_TOKEN_KEY == "" {
		err = errors.New("env REFRESH_TOKEN_KEY must be provided to complete auth flow")
		return
	}

	if config.API_KEY == "" {
		err = errors.New("env API_KEY must be provided to complete auth flow")
	}

	if config.ACCESS_TOKEN_TIME == "" {
		err = errors.New("env ACCESS_TOKEN_TIME must be provided to complete auth flow")
	}

	if config.REFRESH_TOKEN_TIME == "" {
		err = errors.New("env REFRESH_TOKEN_TIME must be provided to complete auth flow")
	}

	if config.WEBSTIE == "" {
		err = errors.New("env WEBSTIE must be provided to complete auth flow")
	}

	if config.DOMAIN == "" {
		err = errors.New("env DOMAIN must be provided to complete auth flow")
	}

	return
}
