package dto

import "time"

type SurveyDTO struct {
	SurveyID           string
	Name               string
	IMO                string
	LoadingDate        time.Time
	LoadingCountry     string
	LoadingPort        string
	DestinationCountry string
	Operation          string
	CargoOnBoard       float64
	Status             string
}
