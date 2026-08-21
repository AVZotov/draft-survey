// Command verify runs CalcDraft/CalcSurvey against a JSON-encoded
// types.Survey and prints the SurveyResult as JSON, for comparing
// calculation output against expected values without going through the
// HTTP/DB stack.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AVZotov/draft-survey/internal/calculation"
	"github.com/AVZotov/draft-survey/internal/types"
)

// verifyOutput adds DraftIndex to each draft total — types.DraftTotal has
// no such field, since the app identifies drafts by slice position instead.
type verifyOutput struct {
	SurveyResult verifySurveyResult `json:"survey_result"`
}

type verifySurveyResult struct {
	CargoOnBoard      float64            `json:"cargo_on_board"`
	CargoDiffDeclared float64            `json:"cargo_diff_declared"`
	DraftTotals       []verifyDraftTotal `json:"draft_totals"`
}

type verifyDraftTotal struct {
	DraftIndex    int               `json:"draft_index"`
	DraftResult   types.DraftResult `json:"draft_result"`
	TrueConstant  float64           `json:"true_constant"`
	CargoFromPrev float64           `json:"cargo_from_prev"`
	ConstantDiff  float64           `json:"constant_diff"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: verify <input.json>")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	var survey types.Survey
	if err := json.Unmarshal(data, &survey); err != nil {
		return fmt.Errorf("unmarshal input: %w", err)
	}

	sr := calculation.CalcSurvey(survey)

	out := verifyOutput{
		SurveyResult: verifySurveyResult{
			CargoOnBoard:      sr.CargoOnBoard,
			CargoDiffDeclared: sr.CargoDiffDeclared,
			DraftTotals:       make([]verifyDraftTotal, len(sr.DraftTotals)),
		},
	}
	for i, dt := range sr.DraftTotals {
		out.SurveyResult.DraftTotals[i] = verifyDraftTotal{
			DraftIndex:    i,
			DraftResult:   dt.DraftResult,
			TrueConstant:  dt.TrueConstant,
			CargoFromPrev: dt.CargoFromPrev,
			ConstantDiff:  dt.ConstantDiff,
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
