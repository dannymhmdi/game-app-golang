package params

import "mymodule/entity"

type ProfileRequest struct {
	Id uint `json:"id"`
}

type ProfileResponse struct {
	RegeneratedToken string      `json:"regenerated_token"`
	User             entity.User `json:"user"`
}
