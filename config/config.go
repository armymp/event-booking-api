package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port int
		Env  string
	}
	Database struct {
		Name string
	}
	JWT struct {
		Secret string
	}
}

var AppConfig Config

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println(".env file not found, relying on environment variables")
    }

    viper.SetDefault("APP_ENV", "development")
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()

    viper.BindEnv("JWT.Secret", "JWT_SECRET")
    viper.BindEnv("Database.Name", "DB_NAME")

    viper.SetConfigName("config." + env)
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    _ = viper.MergeInConfig() // ignore missing YAML

    if err := viper.Unmarshal(&AppConfig); err != nil {
        panic(fmt.Errorf("config unmarshal error: %w", err))
    }

    if AppConfig.JWT.Secret == "" {
        panic("JWT_SECRET is not set")
    }

    fmt.Println("Loaded environment config:", env)
}