package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	LogLevel    string `yaml:"loglevel" env:"DMARCANALYZE_LOG_LEVEL" `
	LogFormat   string `yaml:"logformat" env:"DMARCANALYZE_LOG_FORMAT" `
	LogProgress int    `yaml:"logprogress" env:"DMARCANALYZE_LOG_PROGRESS"`
	Sleep       int    `yaml:"sleep" env:"DMARCANALYZE_SLEEP"`

	IMAP struct {
		Address  string `yaml:"address" env:"DMARCANALYZE_IMAP_SERVER_ADDRESS" `
		Port     string `yaml:"port" env:"DMARCANALYZE_IMAP_SERVER_PORT" `
		Username string `yaml:"username" env:"DMARCANALYZE_IMAP_USERNAME"`
		Password string `yaml:"password" env:"DMARCANALYZE_IMAP_PASSWORD"`
	} `yaml:"imap"`

	Database struct {
		Driver           string `yaml:"driver" env:"DMARCANALYZE_DATABASE_DRIVER" `
		ConnectionString string `yaml:"connectionstring" env:"DMARCANALYZE_DATABASE_CONNECTIONSTRING"`
	} `yaml:"database"`
}

type OffHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

func NewOffHandler() *OffHandler {
	// Optimization: avoid chains of LevelHandlers.
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	return &OffHandler{slog.LevelError.Level(), h}
}

// Never allow logging
func (h *OffHandler) Enabled(_ context.Context, level slog.Level) bool {
	return false
}

// Handle implements Handler.Handle.
func (h *OffHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *OffHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewOffHandler()
}

// WithGroup implements Handler.WithGroup.
func (h *OffHandler) WithGroup(name string) slog.Handler {
	return NewOffHandler()
}

// Handler returns the Handler wrapped by h.
func (h *OffHandler) Handler() slog.Handler {
	return h.handler
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
	if Configuration.LogFormat != "off" {
		slog.Info("log level configured", "level", LogLevel)
	}

	switch Configuration.LogFormat {
	case "off":
		LogLevel.Set(slog.LevelError)
		slog.SetDefault(slog.New(NewOffHandler()))
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: LogLevel,
		})))
	default: // text
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: LogLevel,
		})))
	}

	slog.Debug("configuration", "config", Configuration)
}

func init() {
	ReadConfig()
}
