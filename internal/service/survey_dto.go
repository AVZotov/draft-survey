package service

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type SurveyDTO struct {
	SurveyID           string
	Name               string
	IMO                string
	SurveyDate         string
	SurveyTime         string
	LoadingCountry     string
	LoadingPort        string
	DestinationCountry string
	Operation          string
	CargoOnBoard       float64
	Status             types.SurveyStatus
}

func ToSurveyDTO(s *types.Survey) SurveyDTO {
	return SurveyDTO{
		SurveyID:   s.ID,
		Name:       s.VesselData.Name,
		IMO:        s.VesselData.IMO,
		SurveyDate: s.CreatedAt.Format("02 Jan 2006"),
		SurveyTime: s.CreatedAt.Format("15:04"),
		Operation:  string(s.CargoOperation.Operation),
		Status:     s.Status,
	}
}
