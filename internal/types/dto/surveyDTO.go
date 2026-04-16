package dto

import (
	"time"

	"github.com/AVZotov/draft-survey/internal/types"
)

type SurveyDTO struct {
	SurveyID           string
	Name               string
	IMO                string
	SurveyDate         time.Time
	LoadingCountry     string
	LoadingPort        string
	DestinationCountry string
	Operation          string
	CargoOnBoard       float64
	Status             types.SurveyStatus
}
