package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
}

func NewSlogLogger(level Level) *SlogLogger {
	return &SlogLogger{
		log: slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: toSlogLevel(level),
				},
			),
		),
	}
}

func toSlogLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *SlogLogger) Info(msg string, fields ...any) {
	l.log.Info(msg, fields...)
}

func (l *SlogLogger) Warn(msg string, fields ...any) {
	l.log.Warn(msg, fields...)
}

func (l *SlogLogger) Error(msg string, err error, fields ...any) {
	l.log.Error(msg, append([]any{"error", err}, fields...)...)
}

func (l *SlogLogger) Audit(msg string, surveyID string, fields ...any) {
	l.log.Info(msg, append([]any{"audit", true, "survey_id", surveyID}, fields...)...)
}

func (l *SlogLogger) Debug(msg string, fields ...any) {
	l.log.Debug(msg, fields...)
}
