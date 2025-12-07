package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
		Env  string
	}
	Database struct {
		Name string
	}
}

var AppConfig Config

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // default
	}

	// Load the config file based on env
	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look in current folder

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("config file error: %w", err))
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("config unmarshal error: %w", err))
	}

	fmt.Println("Loaded environment config:", env)
}