package handler

import (
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/validation"
)

type Handler struct {
	userRepository         storage.UserRepository
	surveyRepository       storage.SurveyRepository
	dictionariesRepository storage.DictionariesRepository
	validator              validation.Validator
	appVersion             string
}

func New(
	userRepository storage.UserRepository,
	surveyRepository storage.SurveyRepository,
	dictionariesRepository storage.DictionariesRepository,
	validator validation.Validator,
	version string,
) *Handler {
	return &Handler{
		userRepository:         userRepository,
		surveyRepository:       surveyRepository,
		dictionariesRepository: dictionariesRepository,
		validator:              validator,
		appVersion:             version,
	}
}
