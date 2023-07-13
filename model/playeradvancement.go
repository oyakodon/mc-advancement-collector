package model

import "time"

type PlayerAdvancementSummary struct {
	Advancements map[string]*PlayerAdvancement
	Progress     AdvancementProgress
	Updated      time.Time
	Cached       time.Time
}
type PlayerAdvancement struct {
	Key      string                   `json:"key"`
	Parent   string                   `json:"parent"`
	Display  PlayerAdvancementDisplay `json:"display"`
	Type     AdvancementType          `json:"type"`
	Hidden   bool                     `json:"hidden"`
	Done     bool                     `json:"done"`
	Criteria map[string]*time.Time    `json:"criteria"`
	Progress AdvancementProgress      `json:"progress"`
}

type PlayerAdvancementDisplay struct {
	Title       string                       `json:"title"`
	Description string                       `json:"desciption"`
	Icon        PlayerAdvancementDisplayIcon `json:"icon"`
}

type PlayerAdvancementDisplayIcon struct {
	Url       string `json:"url"`
	InvSprite bool   `json:"invsprite"`
	PosX      *int   `json:"posx,omitempty"`
	PosY      *int   `json:"posy,omitempty"`
}
