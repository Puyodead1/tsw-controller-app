package logger

import (
	"log/slog"
	"os"
)

var loghandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})

var Logger = slog.New(loghandler)
