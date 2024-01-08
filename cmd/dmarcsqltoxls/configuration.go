package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	LogLevel    string `yaml:"loglevel" env:"DMARCSQLTOXLS_LOG_LEVEL" `
	LogFormat   string `yaml:"logformat" env:"DMARCSQLTOXLS_LOG_FORMAT" `
	LogProgress int    `yaml:"logprogress" env:"DMARCSQLTOXLS_LOG_PROGRESS"`

	Database struct {
		Driver           string `yaml:"driver" env:"DMARCSQLTOXLS_DATABASE_DRIVER" `
		ConnectionString string `yaml:"connectionstring" env:"DMARCSQLTOXLS_DATABASE_CONNECTIONSTRING"`
	} `yaml:"database"`

	XLS struct {
		Template string `yaml:"template" env:"DMARCSQLTOXLS_XLS_TEMPLATE"`
		Output   string `yaml:"output" env:"DMARCSQLTOXLS_XLS_OUTPUT"`
		Style    struct {
			Datasheet struct {
				DateFormat          string `yaml:"dateFormat"`
				RowABackgroundLight string `yaml:"rowabackgroundlight"`
				RowABackgroundDark  string `yaml:"rowabackgrounddark"`
				RowAFailFontColor   string `yaml:"rowafailfontcolor"`
				RowBBackgroundLight string `yaml:"rowbbackgroundlight"`
				RowBBackgroundDark  string `yaml:"rowbbackgrounddark"`
				RowBFailFontColor   string `yaml:"rowbfailfontcolor"`
			} `yaml:"datasheet"`
		} `yaml:"style"`
	} `yaml:"xls"`
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
