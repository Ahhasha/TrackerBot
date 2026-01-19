package config

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidConfig = errors.New("invalid config")

type Config struct {
	TelegramToken string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

func Load() (Config, error) {
	cfg := Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string
	if c.TelegramToken == "" {
		missing = append(missing, "TELEGRAM_BOT_TOKEN")
	}
	if c.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing env param: %v", missing)
	}
	return nil
}
