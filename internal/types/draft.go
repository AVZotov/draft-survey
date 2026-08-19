package types

// NewDraft returns a Draft initialized with the fields required by the
// draft initialization invariant (see internal/CLAUDE.md): HydrostaticRows
// and MTCRows must always have length 2, or hydrostatics.templ panics on
// an unconditional rows[0]/rows[1] access. PP directions default to the
// standard bow/stern convention (Fwd/Mid measured aft, Aft measured
// forward).
//
// Type and Status are deliberately left unset — callers set them
// explicitly, since Survey.Create() and DraftService.Add() each need a
// different Type/Status combination for the draft they're creating.
func NewDraft() Draft {
	return Draft{
		HydrostaticRows: make([]HydrostaticRow, 2),
		MTCRows:         make([]MTCRow, 2),
		PPFwdDirection:  "A",
		PPMidDirection:  "A",
		PPAftDirection:  "F",
	}
}
