package configs

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func ConfigureLogger(level slog.Leveler) {
	logOpts := new(tint.Options)
	logOpts.Level = level
	logOpts.AddSource = false
	logOpts.NoColor = false
	logOpts.TimeFormat = ""
	handler := tint.NewHandler(os.Stdout, logOpts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
