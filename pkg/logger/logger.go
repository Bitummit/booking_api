package logger

import (
	"log/slog"
	"os"
)


func NewLogger() *slog.Logger{
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}

func Err(err error) slog.Attr{
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}