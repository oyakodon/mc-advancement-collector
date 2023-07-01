package model

type AdvancementType string

const (
	Task      AdvancementType = "task"
	Goal      AdvancementType = "goal"
	Challenge AdvancementType = "challenge"
)

type CalculateType string

const (
	CalculateOneOf CalculateType = "oneof"
	CalculateAllOf CalculateType = "allof"
)

type AdvancementFilterCondition string

const (
	ConditionAll      AdvancementFilterCondition = "all"
	ConditionDone     AdvancementFilterCondition = "done"
	ConditionProgress AdvancementFilterCondition = "progress"
)
