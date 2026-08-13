package service

type NotificationKind string

const (
	KindSuccess NotificationKind = "success"
	KindWarning NotificationKind = "warn"
	KindError   NotificationKind = "error"
	KindInfo    NotificationKind = "info"
)

type Notification struct {
	Kind    NotificationKind `json:"kind"`
	Header  string           `json:"header"`
	Message string           `json:"message"`
}

type Outcome struct {
	Toast    *Notification
	Alert    *Notification
	Redirect string
}
