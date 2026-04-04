package constants

// Alert thresholds for draft survey warnings
// These are initial values — to be calibrated with surveyors and auditors

// Trim thresholds, meters
const (
	// TrimHeadWarning — trim by head warning threshold
	// Non-standard condition, calculations still reliable
	TrimHeadWarning = 0.5

	// TrimHeadDanger — trim by head danger threshold
	// Tank volume calculations become unreliable, LOP recommended
	TrimHeadDanger = 1.0

	// TrimEvenKeel — even keel tolerance
	// Trim within this range considered even keel
	TrimEvenKeel = 0.1
)

// List thresholds, degrees
const (
	// ListWarning — list warning threshold per UNECE/DSGear standard
	// LOP to be issued
	ListWarning = 0.5

	// ListDanger — list danger threshold
	// Beyond calibration table range, tank calculations unreliable
	ListDanger = 1.0
)

// Deflection thresholds, cm
const (
	// DeflectionWarning — deflection warning threshold
	DeflectionWarning = 1.0

	// DeflectionDanger — deflection danger threshold
	DeflectionDanger = 5.0
)
