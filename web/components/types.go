package components

import (
	"github.com/AVZotov/draft-survey/internal/types"
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
}

type CountrySelectProps struct {
	Countries []types.Country
	Selected  string
	Name      string
}

type PortSelectProps struct {
	Ports    []types.Port
	Selected string
	Name     string
}

type SurveyPageProps struct {
	Survey types.Survey
}

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
