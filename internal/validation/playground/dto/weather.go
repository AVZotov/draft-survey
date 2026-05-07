package dto

import "github.com/AVZotov/draft-survey/internal/types"

type SeaCondition struct {
	Type string `validate:"oneof=wave ice"`
	Wave string `validate:"required_if=Type wave"`
	Ice  string `validate:"required_if=Type ice"`
}

func NewSeaConditionDTO(sc *types.SeaCondition) *SeaCondition {
	if sc == nil {
		return nil
	}
	return &SeaCondition{
		Type: string(sc.Type),
		Wave: string(sc.Wave),
		Ice:  string(sc.Ice),
	}
}
