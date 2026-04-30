package dto

import "github.com/AVZotov/draft-survey/internal/types"

type AlpineDTO struct {
	Survey        types.Survey       `json:"survey"`
	SurveyResults types.SurveyResult `json:"survey_results"`
}
