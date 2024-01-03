package main

import (
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	LogLevel  string `yaml:"loglevel" env:"DMARCANALYZE_LOG_LEVEL" env-default:"debug"`
	LogFormat string `yaml:"logformat" env:"DMARCANALYZE_LOG_FORMAT" env-default:"text"`

	IMAP struct {
		Address  string `yaml:"address" env:"DMARCANALYZE_IMAP_SERVER_ADDRESS" env-default:""`
		Port     string `yaml:"port" env:"DMARCANALYZE_IMAP_SERVER_PORT" env-default:""`
		Username string `yaml:"username" env:"DMARCANALYZE_IMAP_USERNAME" env-default:""`
		Password string `yaml:"password" env:"DMARCANALYZE_IMAP_PASSWORD" env-default:""`
	} `yaml:"imap"`

	Database struct {
		Driver string `yaml:"driver" env:"DMARCANALYZE_DATABASE_DRIVER" env-default:"sqlite3"`
		DSN    string `yaml:"dsn" env:"DMARCANALYZE_DATABASE_DSN" env-default:"dmarc.db"`
	} `yaml:"database"`
}

var (
	LogLevel      = new(slog.LevelVar)
	Configuration ConfigDatabase
)

func ReadConfig() {
	LogLevel.Set(slog.LevelDebug)
	err := cleanenv.ReadConfig("config.yml", &Configuration)
	if err != nil {
		slog.Error("error reading configuration", "error", err)
	}

	switch Configuration.LogLevel {
	case "info":
		LogLevel.Set(slog.LevelInfo)
	case "warn":
		LogLevel.Set(slog.LevelWarn)
	case "error":
		LogLevel.Set(slog.LevelError)
	default:
		LogLevel.Set(slog.LevelDebug)
	}
	slog.Info("log level configured", "level", LogLevel)
	if Configuration.LogFormat == "json" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: LogLevel,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: LogLevel,
		})))
	}
	slog.Debug("configuration", "config", Configuration)
}

func init() {
	ReadConfig()
}
