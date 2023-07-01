package responses

import "com.oykdn.mc-advancement-collector/config"

type AdvancementAssetsResponse struct {
	Background map[string]AdvancementAssetsBackground `json:"background"`
}

type AdvancementAssetsBackground struct {
	Incomplete string `json:"incomplete"`
	Completed  string `json:"completed"`
}

func ConvertToAdvancementAssetsResponse(conf map[string]config.AppConfigAssetBackground) *AdvancementAssetsResponse {
	background := make(map[string]AdvancementAssetsBackground)
	for k, v := range conf {
		background[k] = AdvancementAssetsBackground{
			Incomplete: v.Incomplete,
			Completed:  v.Completed,
		}
	}

	return &AdvancementAssetsResponse{
		Background: background,
	}
}
