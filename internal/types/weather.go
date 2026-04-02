package types

type SeaConditionType string

const (
	SeaConditionTypeWave SeaConditionType = "wave"
	SeaConditionTypeIce  SeaConditionType = "ice"
)

type WaveCondition string

const (
	WaveConditionCalm     WaveCondition = "under 0.1m"
	WaveConditionSmooth   WaveCondition = "0.1-0.5m"
	WaveConditionSlight   WaveCondition = "0.5-1.25m"
	WaveConditionModerate WaveCondition = "1.25-2.5m"
	WaveConditionRough    WaveCondition = "2.5-4.0m"
)

type IceCondition string

const (
	IceConditionUnder005 IceCondition = "under 0.05m around"
	IceCondition005To010 IceCondition = "0.05-0.1m around"
	IceCondition010To015 IceCondition = "0.1-0.15m around"
	IceCondition015To020 IceCondition = "0.15-0.2m around"
	IceCondition020To030 IceCondition = "0.2-0.3m around"
	IceCondition030To040 IceCondition = "0.3-0.4m around"
	IceCondition040To060 IceCondition = "0.4-0.6m around"
	IceConditionOver060  IceCondition = "over 0.6m around"
)

type SeaCondition struct {
	Type SeaConditionType `json:"type"`
	Wave WaveCondition    `json:"wave"`
	Ice  IceCondition     `json:"ice"`
}
