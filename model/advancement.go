package model

type MinecraftAdvancement struct {
	Criteria map[string]string `json:"criteria"`
	Done     bool              `json:"done"`
}

type AdvancementProgress struct {
	Total      int     `json:"total"`
	Done       int     `json:"done"`
	Percentage float64 `json:"percentage"`
}
