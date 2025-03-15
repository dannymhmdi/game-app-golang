package params

import "mymodule/entity"

type MatchMakingRequest struct{}

type MatchMakingResponse struct {
	MatchedUsers []entity.MatchedPlayers
}
