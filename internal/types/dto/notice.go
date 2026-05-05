package dto

type NoticeType string

const (
	NoticeError NoticeType = "error"
	NoticeWarn  NoticeType = "warn"
	NoticeInfo  NoticeType = "info"
)

type NoticeItem struct {
	Type    NoticeType
	Field   string
	Message string
}
