package requests

import "com.oykdn.mc-advancement-collector/model"

type PlayerAdvancementRequest struct {
	PlayerId  string                           `uri:"id" binding:"required,uuid"`
	Condition model.AdvancementFilterCondition `form:"condition"`
}
