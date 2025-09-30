package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AccessTokenExp time.Duration
}

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ACCESS_TOKEN_EXP", "15m")

	cfg := &Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
	}

	exp := viper.GetString("ACCESS_TOKEN_EXP")
	dur, _ := time.ParseDuration(exp)
	cfg.AccessTokenExp = dur

	return cfg, nil
}
