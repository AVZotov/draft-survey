package chi

import (
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
)

type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	Audit(msg string, surveyID string, fields ...any)
}

type Handler struct {
	userRepository         storage.UserRepository
	surveyRepository       storage.SurveyRepository
	surveyQueryRepository  storage.SurveyQueryRepository
	dictionariesRepository storage.DictionariesRepository
	validator              validation.Validator
	logger                 Logger
	appVersion             string
}

func New(
	userRepository storage.UserRepository,
	surveyRepository storage.SurveyRepository,
	surveyQueryRepository storage.SurveyQueryRepository,
	dictionariesRepository storage.DictionariesRepository,
	validator validation.Validator,
	logger Logger,
	version string,
) *Handler {
	return &Handler{
		userRepository:         userRepository,
		surveyRepository:       surveyRepository,
		surveyQueryRepository:  surveyQueryRepository,
		dictionariesRepository: dictionariesRepository,
		validator:              validator,
		logger:                 logger,
		appVersion:             version,
	}
}
