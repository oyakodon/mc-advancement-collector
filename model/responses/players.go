package responses

import "com.oykdn.mc-advancement-collector/model"

type PlayersResponse struct {
	Players []model.PlayerProfile `json:"players"`
}
