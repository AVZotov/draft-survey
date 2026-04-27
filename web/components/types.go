package components

import (
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/types/dto"
)

type LayoutProps struct {
	Title           string
	MetaDescription string
	User            *types.User
	Survey          *types.Survey
	Results         *[]types.DraftResult
	ExtraCSS        []string
	ExtraJS         []string
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
	Survey    types.Survey
	Countries []types.Country
}

type CountrySelectProps struct {
	Countries     []types.Country
	SelectedCode  string // country code e.g. "RU"
	SelectedName  string // country name e.g. "Russia"
	CountryCodeID string // hidden input name for code e.g. "country_code"
	CountryNameID string // hidden input name for name e.g. "country_name"
}

type PortSelectProps struct {
	Ports    []types.Port
	Selected string
	Name     string
}

type SurveyListProps struct {
	Total      int
	Complete   int
	InProgress int
	Draft      int
	AllDTO     []dto.SurveyDTO
	WeekDTO    []dto.SurveyDTO
	MonthDTO   []dto.SurveyDTO
}

type ResultProps struct {
	Survey    dto.SurveyDTO
	Alpine    string
	Lastindex int
}

//To be reviewed props

type Version struct {
	Version string
}

type BannerType string

const (
	Info  BannerType = "info"
	Warn  BannerType = "warn"
	Error BannerType = "error"
)

type BannerProps struct {
	Type    BannerType
	Header  string
	Message string
	Details string
}

type ModalLevel string

const (
	ModalInfo    ModalLevel = "Information"
	ModalWarning ModalLevel = "Warning"
	ModalError   ModalLevel = "Danger"
)

type ModalProps struct {
	Level      ModalLevel
	DialogID   string
	Title      string
	Message    string
	ConfirmBtn string
	CancelBtn  string // Do not rendered if prop is empty
}
