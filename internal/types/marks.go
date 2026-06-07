package types

type ReadingMethod string

const (
	ReadingMethodDirect ReadingMethod = "direct"
	ReadingMethodCamera ReadingMethod = "camera"
)

type Mark struct {
	Value  *float64      `json:"value"`
	Method ReadingMethod `json:"method"`
}

func (rm ReadingMethod) Short() string {
	if len(rm) < 3 {
		return string(rm)
	}
	return string(rm[:3])
}

type Marks struct {
	FwdPort      Mark `json:"fwd_port"`
	FwdStarboard Mark `json:"fwd_stbd"`
	MidPort      Mark `json:"mid_port"`
	MidStarboard Mark `json:"mid_stbd"`
	AftPort      Mark `json:"aft_port"`
	AftStarboard Mark `json:"aft_stbd"`
}
