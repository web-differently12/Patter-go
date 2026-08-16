package core

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
	return Logger
}

func GetLogger() *slog.Logger {
	if Logger == nil {
		return InitLogger()
	}
	return Logger
}
