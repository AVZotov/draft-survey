package handler

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

type HandlerChi struct {
	userRepository         storage.UserRepository
	surveyRepository       storage.SurveyRepository
	surveyQueryRepository  storage.SurveyQueryRepository
	dictionariesRepository storage.DictionariesRepository
	validator              validation.Validator
	logger                 Logger
	appVersion             string
}

func NewChi(
	userRepository storage.UserRepository,
	surveyRepository storage.SurveyRepository,
	surveyQueryRepository storage.SurveyQueryRepository,
	dictionariesRepository storage.DictionariesRepository,
	validator validation.Validator,
	logger Logger,
	version string,
) *HandlerChi {
	return &HandlerChi{
		userRepository:         userRepository,
		surveyRepository:       surveyRepository,
		surveyQueryRepository:  surveyQueryRepository,
		dictionariesRepository: dictionariesRepository,
		validator:              validator,
		logger:                 logger,
		appVersion:             version,
	}
}
