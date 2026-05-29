package chi

import (
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/service"
)

type Handler struct {
	services   *service.Services
	logger     logger.Logger
	appVersion string
}

func New(services *service.Services, logger logger.Logger, version string) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		appVersion: version,
	}
}
