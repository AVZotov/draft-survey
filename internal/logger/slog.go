package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{
		log: slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		),
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
