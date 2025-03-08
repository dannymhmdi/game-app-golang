package params

type PresenseRequest struct {
	UserId    uint  `json:"user_id"`
	Timestamp int64 `json:"timestamp"`
}

type PresenseResponse struct {
	Message string
}
