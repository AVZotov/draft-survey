package chi

import (
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/service"
	"github.com/AVZotov/draft-survey/internal/sse"
)

type Handler struct {
	services   *service.Services
	logger     logger.Logger
	appVersion string
	broker     *sse.Broker
}

func New(services *service.Services, logger logger.Logger, version string, broker *sse.Broker) *Handler {
	return &Handler{
		services:   services,
		logger:     logger,
		appVersion: version,
		broker:     broker,
	}
}
