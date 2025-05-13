package params

import "mymodule/entity"

type MatchStoreRequest struct {
	GameInfo entity.Game `json:"game_info"`
}

type MatchStoreResponse struct {
	Message string `json:"message"`
	GameId  uint   `json:"game_id"`
}
