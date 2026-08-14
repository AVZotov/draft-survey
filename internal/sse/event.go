package sse

type EventType string

const (
	EventToast       EventType = "toast"
	EventAlert       EventType = "alert"
	EventSurveyStats EventType = "survey-stats"
	EventDraftCalc   EventType = "draft-calc"
)

type Event struct {
	Type EventType
	Data string
}
