package format

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/vessel"
)

func Mark(p *float64) string {
	if p == nil {
		return "—"
	}
	return fmt.Sprintf("%.3f", *p)
}

func FloatOrEmpty(v *float64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%g", *v)
}

func IntOrEmpty(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

func IntToStr(i int) string {
	if i == 0 {
		return "-"
	}
	return strconv.Itoa(i)
}

func Weight(v float64) string {
	return fmt.Sprintf("%.3f", v)
}

func WeightFormatted(v float64) string {
	negative := v < 0
	v = math.Abs(v)
	v = math.Round(v*1000) / 1000
	sep := 1000
	var str string
	whole := int(math.Trunc(v))
	var rem int
	fractional := int(math.Trunc((v - float64(whole)) * 1000))
	var data []int

	for {
		rem = whole % sep
		data = slices.Insert(data, 0, rem)

		if whole/sep == 0 {
			for i, v := range data {
				if i == 0 {
					str += strconv.Itoa(v)
					continue
				}
				str += fmt.Sprintf(" %03d", v)
			}
			break
		}
		whole = whole / sep
	}
	str += fmt.Sprintf(".%03d", fractional)
	if negative {
		str = fmt.Sprintf("-%s", str)
	}
	return str
}

func Draft(v float64) string {
	return fmt.Sprintf("%.3f m", v)
}

func LBPCorrection(correction vessel.CorrectionMethod) string {
	lbpCorrections := vessel.CorrectionMethodFullLBP
	if correction == vessel.CorrectionMethodHalfLBP {
		lbpCorrections = vessel.CorrectionMethodHalfLBP
	}

	return string(lbpCorrections)
}

func MMCShortFormula(vesselType vessel.VesselType) string {
	mmcCorrections := "(6m/8)"
	if vesselType == vessel.VesselTypeRiver {
		mmcCorrections = "(4m/6)"
	}
	if vesselType == vessel.VesselTypeBarge {
		mmcCorrections = "(14m/20)"
	}

	return mmcCorrections
}

// TrimDirection returns the trim direction as a string.
// Positive trim = stern deeper = "Aft", negative = bow deeper = "Fwd"
func TrimDirection(trim float64) string {
	switch {
	case trim > 0:
		return "Aft"
	case trim < 0:
		return "Fwd"
	default:
		return "Even keel"
	}
}

// ListDirection returns the list direction as a string.
// Positive list = port side deeper = "Port", negative = starboard deeper = "Stbd"
func ListDirection(list float64) string {
	switch {
	case list > 0:
		return "Port"
	case list < 0:
		return "Stbd"
	default:
		return ""
	}
}

func TankNameID(prefix string, tank types.Tank) string {
	return fmt.Sprintf("%s-%s", prefix, tank.ID)
}
