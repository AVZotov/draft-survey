package sse

type EventType string

const (
	EventToast       EventType = "toast"
	EventAlert       EventType = "alert"
	EventSurveyStats EventType = "survey-stats"
)

type Event struct {
	Type EventType
	Data string
}
