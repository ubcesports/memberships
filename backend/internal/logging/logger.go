package logging

import (
	"log/slog"
	"os"
)

func Configure() {
	level := slog.LevelInfo // levels go from debug -> info -> warn -> error

	// allow debug logs for development only. they don't need to clutter the prod logs
	if os.Getenv("APP_ENV") == "development" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}).WithAttrs([]slog.Attr{
		slog.String("environment", os.Getenv("APP_ENV")),
	})
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
