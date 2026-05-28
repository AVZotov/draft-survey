package chi

import (
	"github.com/AVZotov/draft-survey/internal/service"
)

type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	Audit(msg string, surveyID string, fields ...any)
}

type Handler struct {
	services   *service.Services
	logger     Logger
	appVersion string
}

func New(services *service.Services, logger Logger, version string) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		appVersion: version,
	}
}
