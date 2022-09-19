package config

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Debug  bool `yaml:"debug"`
	Server struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Https struct {
			Enabled bool   `yaml:"enabled"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"https"`
		RateLimit struct {
			Enabled  bool          `yaml:"enabled"`
			Tokens   uint64        `yaml:"tokens"`
			Interval time.Duration `yaml:"interval"`
		} `yaml:"rate-limit"`
	} `yaml:"server"`
	Datasource struct {
		Host          string `yaml:"host"`
		Port          string `yaml:"port"`
		User          string `yaml:"user"`
		Password      string `yaml:"password"`
		Database      string `yaml:"database"`
		Ssl           bool   `yaml:"ssl"`
		Timezone      string `yaml:"timezone"`
		AutoMigration bool   `yaml:"auto-migration"`
	} `yaml:"datasource"`
	Telegram struct {
		Token  string `yaml:"token"`
		ChatId int64  `yaml:"chat"`
	} `yaml:"telegram"`
}

func newConfig(configPath string) Config {
	// Create config structure
	config := defaultConfig()
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.WithError(err).Warn("Config: reading error")
		return config
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.WithError(err).Warn("Config: parsing error")
	}
	return config
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		message := fmt.Sprintf("'%s' is a directory, not a normal file", path)
		err := errors.New(message)
		log.Warn(message, err)
		return err
	}
	return nil
}

func parseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}

func FromFile() Config {
	path, err := parseFlags()
	if err != nil {
		log.WithError(err).Warn("Config: completed with error")
		return defaultConfig()
	}
	config := newConfig(path)
	log.Info("Config: completed")
	return config
}

func defaultConfig() Config {
	return Config{
		Debug: os.Getenv("DEBUG") == "true",
		Server: struct {
			Host  string `yaml:"host"`
			Port  string `yaml:"port"`
			Https struct {
				Enabled bool   `yaml:"enabled"`
				Cert    string `yaml:"cert"`
				Key     string `yaml:"key"`
			} `yaml:"https"`
			RateLimit struct {
				Enabled  bool          `yaml:"enabled"`
				Tokens   uint64        `yaml:"tokens"`
				Interval time.Duration `yaml:"interval"`
			} `yaml:"rate-limit"`
		}{
			Host: "127.0.0.1",
			Port: "8080",
		},
		Datasource: struct {
			Host          string `yaml:"host"`
			Port          string `yaml:"port"`
			User          string `yaml:"user"`
			Password      string `yaml:"password"`
			Database      string `yaml:"database"`
			Ssl           bool   `yaml:"ssl"`
			Timezone      string `yaml:"timezone"`
			AutoMigration bool   `yaml:"auto-migration"`
		}{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
			Timezone: "UTC",
		},
	}
}
