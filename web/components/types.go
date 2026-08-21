package components

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type LayoutProps struct {
	Title           string
	MetaDescription string
	User            *types.User
	ExtraCSS        []string
	ExtraJS         []string
	AppVersion      string
}
type TanksPageProps struct {
	Survey     types.Survey
	DraftIndex int
	Trim       *float64
	TrimDir    string
	List       *float64
	ListDir    string
}

type DraftPageProps struct {
	Survey types.Survey
}

type SeaConditionProps struct {
	SeaCondition types.SeaCondition
	DraftIndex   *int
	SurveyID     string
}

type SurveyPageProps struct {
	Survey types.Survey
}

type PortSelectProps struct {
	Ports    []types.Port
	Selected string
	Name     string
}
type ResultProps struct {
	Survey       types.Survey
	SurveyResult types.SurveyResult
	FirstIndex   int
	LastIndex    int
}
