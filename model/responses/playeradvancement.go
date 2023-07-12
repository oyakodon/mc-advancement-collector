package responses

import (
	"time"

	"com.oykdn.mc-advancement-collector/model"
)

type PlayerAdvancementResponse struct {
	Advancements []*model.PlayerAdvancement `json:"advancements"`
	Progress     model.AdvancementProgress  `json:"progress"`
	Updated      time.Time                  `json:"updated"`
	Cached       time.Time                  `json:"cached"`
}
