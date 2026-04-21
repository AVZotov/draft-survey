package types

import "time"

type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"` // "warning", "override", "error"
	Message   string    `json:"message"`
	UserID    string    `json:"user_id"`
}
