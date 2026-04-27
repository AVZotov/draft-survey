package dto

import "github.com/AVZotov/draft-survey/internal/types"

type AlpineDTO struct {
	Drafts      []types.Draft      `json:"drafts"`
	DraftTotals []types.DraftTotal `json:"draft_totals"`
}
