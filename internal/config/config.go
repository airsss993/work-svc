package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/airsss993/work-svc/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server    Server
		LDAP      LDAPConfig
		GitBucket GitBucketConfig
		App       App
	}
	Server struct {
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}

	App struct {
		Test   bool
		WebURL string
	}
	LDAPConfig struct {
		URL string
	}

	GitBucketConfig struct {
		URL    string
		APIKey string
	}
)

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using system environment variables")
	}

	if err := parseConfigFile("./config"); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := setFromEnv(&cfg); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to set environment variables: %w", err)
	}

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) error {
	cfg.LDAP.URL = os.Getenv("LDAP_URL")
	cfg.GitBucket.URL = os.Getenv("GITBUCKET_URL")
	cfg.GitBucket.APIKey = os.Getenv("GITBUCKET_API_KEY")
	cfg.App.WebURL = os.Getenv("WEB_URL")

	if cfg.LDAP.URL == "" {
		return errors.New("LDAP_URL environment variable is required")
	}

	if cfg.GitBucket.URL == "" {
		return errors.New("GITBUCKET_URL environment variable is required")
	}
	if cfg.GitBucket.APIKey == "" {
		return errors.New("GITBUCKET_API_KEY environment variable is required")
	}

	if cfg.App.WebURL == "" {
		return errors.New("WEB_URL environment variable is required")
	}

	return nil
}
