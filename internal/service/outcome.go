package service

type NotificationKind string

const (
	KindSuccess NotificationKind = "success"
	KindWarning NotificationKind = "warning"
	KindError   NotificationKind = "error"
	KindInfo    NotificationKind = "info"
)

type Notification struct {
	Kind    NotificationKind
	Header  string
	Message string
}

type Outcome struct {
	Toast    *Notification
	Alert    *Notification
	Redirect string
}
