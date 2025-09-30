package config

import "time"

type Config struct {
	DBPath         string
	JWTSecret      string
	AccessTokenExp time.Duration
	SMTPHost       string
	SMTPPort       int
	SMTPEmail      string
	SMTPPassword   string
	Port           string
}

func Load() (*Config, error) {
	cfg := &Config{
		DBPath:         "notebooq.db",
		JWTSecret:      "your-secret-key",
		AccessTokenExp: time.Hour * 24,
		SMTPHost:       "smtp.gmail.com",
		SMTPPort:       587,
		Port:           "8080",
	}

	return cfg, nil
}
