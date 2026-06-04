package logger

type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	Audit(msg string, surveyID string, fields ...any)
	Debug(msg string, fields ...any)
}
