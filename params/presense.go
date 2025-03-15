package params

import "mymodule/entity"

type PresenseRequest struct {
	UserId    uint  `json:"user_id"`
	Timestamp int64 `json:"timestamp"`
}

type PresenseResponse struct {
	Message string
}

type GetPresenceRequest struct {
	UserIDs []uint
}

type GetPresenceResponse struct {
	OnlinePlayers []entity.OnlinePlayer
}

