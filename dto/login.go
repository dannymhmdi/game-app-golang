package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Token   Token  `json:"token"`
}
