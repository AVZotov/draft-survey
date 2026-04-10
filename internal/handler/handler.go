package handler

import (
	"github.com/AVZotov/draft-survey/internal/storage"
)

type Handler struct {
	userRepository         storage.UserRepository
	surveyRepository       storage.SurveyRepository
	dictionariesRepository storage.DictionariesRepository
}

func New(userRepository storage.UserRepository, surveyRepository storage.SurveyRepository, dictionariesRepository storage.DictionariesRepository) *Handler {
	return &Handler{
		userRepository:         userRepository,
		surveyRepository:       surveyRepository,
		dictionariesRepository: dictionariesRepository,
	}
}
