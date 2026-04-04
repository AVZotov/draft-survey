package dictionaries

import "github.com/AVZotov/draft-survey/internal/types"

var WaveConditions = []string{
	string(types.WaveConditionCalm),
	string(types.WaveConditionSmooth),
	string(types.WaveConditionSlight),
	string(types.WaveConditionModerate),
	string(types.WaveConditionRough),
}

var IceConditions = []string{
	string(types.IceConditionUnder005),
	string(types.IceCondition005To010),
	string(types.IceCondition010To015),
	string(types.IceCondition015To020),
	string(types.IceCondition020To030),
	string(types.IceCondition030To040),
	string(types.IceCondition040To060),
	string(types.IceConditionOver060),
}
