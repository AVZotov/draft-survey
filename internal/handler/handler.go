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
}

func New(
	userRepository storage.UserRepository,
	surveyRepository storage.SurveyRepository,
	dictionariesRepository storage.DictionariesRepository,
	validator validation.Validator,
) *Handler {
	return &Handler{
		userRepository:         userRepository,
		surveyRepository:       surveyRepository,
		dictionariesRepository: dictionariesRepository,
		validator:              validator,
	}
}
