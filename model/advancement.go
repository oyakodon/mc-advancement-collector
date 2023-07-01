package model

import "time"

type MinecraftAdvancement struct {
	Criteria map[string]string `json:"criteria"`
	Done     bool              `json:"done"`
}

type AdvancementProgress struct {
	Total      int     `json:"total"`
	Done       int     `json:"done"`
	Percentage float64 `json:"percentage"`
}

type MinecraftAdvancementSummary struct {
	Advancements map[string]MinecraftAdvancement
	UpdatedAt    time.Time
}
