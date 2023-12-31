package model

type AdvancementType string

const (
	Task      AdvancementType = "task"
	Goal      AdvancementType = "goal"
	Challenge AdvancementType = "challenge"
)

type MetricsType string

const (
	MetricsAnyOf MetricsType = "anyof"
	MetricsAllOf MetricsType = "allof"
)

type AdvancementFilterCondition string

const (
	ConditionAll      AdvancementFilterCondition = "all"
	ConditionDone     AdvancementFilterCondition = "done"
	ConditionProgress AdvancementFilterCondition = "progress"
)
