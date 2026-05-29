package logger

type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	Audit(msg string, surveyID string, fields ...any)
}

type Noop struct{}

func (l *Noop) Info(msg string, fields ...any)                   {}
func (l *Noop) Warn(msg string, fields ...any)                   {}
func (l *Noop) Error(msg string, err error, fields ...any)       {}
func (l *Noop) Audit(msg string, surveyID string, fields ...any) {}
