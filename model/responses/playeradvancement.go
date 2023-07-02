package responses

import (
	"time"

	"com.oykdn.mc-advancement-collector/model"
)

type PlayerAdvancementResponse struct {
	Advancements map[string]*PlayerAdvancement `json:"advancements"`
	Progress     model.AdvancementProgress     `json:"progress"`
	Updated      time.Time                     `json:"updated"`
	Cached       time.Time                     `json:"cached"`
}

type PlayerAdvancement struct {
	Parent   string                    `json:"parent"`
	Display  PlayerAdvancementDisplay  `json:"display"`
	Type     model.AdvancementType     `json:"type"`
	Hidden   bool                      `json:"hidden"`
	Done     bool                      `json:"done"`
	Criteria map[string]*time.Time     `json:"criteria"`
	Progress model.AdvancementProgress `json:"progress"`
}

type PlayerAdvancementDisplay struct {
	Title       string                       `json:"title"`
	Description string                       `json:"desciption"`
	Icon        PlayerAdvancementDisplayIcon `json:"icon"`
}

type PlayerAdvancementDisplayIcon struct {
	Url       string `json:"url"`
	InvSprite bool   `json:"invsprite"`
	PosX      int    `json:"posx,omitempty"`
	PosY      int    `json:"posy,omitempty"`
}
