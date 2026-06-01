package sse

type EventType string

const (
	EventToast EventType = "toast"
	EventAlert EventType = "alert"
)

type Event struct {
	Type EventType
	Data string
}
