package configs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

const ansiGreen = "\033[32m"
const ansiGray = "\033[90m"
const ansiReset = "\033[0m"

type CustomTextHandler struct {
	level slog.Leveler
}

func NewCustomTextHandler(level slog.Leveler) slog.Handler {
	return &CustomTextHandler{level: level}
}

func (h *CustomTextHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

func (h *CustomTextHandler) Handle(_ context.Context, r slog.Record) error {
	levelStr := h.formatLevel(r.Level)
	fmt.Fprintf(
		os.Stdout,
		ansiGray+"["+ansiGreen+"%s"+ansiReset+ansiGray+"]"+ansiReset+" %s\n",
		levelStr,
		r.Message,
	)
	return nil
}

func (h *CustomTextHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *CustomTextHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *CustomTextHandler) formatLevel(lvl slog.Level) string {
	switch lvl {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return "LOG"
	}
}

func ConfigureLogger(level slog.Leveler) {
	handler := NewCustomTextHandler(level)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
