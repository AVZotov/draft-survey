package validation

import (
	"github.com/AVZotov/draft-survey/internal/types"
)

type DraftValidator interface {
	Validate(draft types.Draft, result types.DraftResult) []types.AuditEvent
}

type NoopDraftValidator struct{}

func (v *NoopDraftValidator) Validate(draft types.Draft, result types.DraftResult) []types.AuditEvent {
	return nil
}

func NewDraftValidator() DraftValidator {
	return &NoopDraftValidator{}
}
