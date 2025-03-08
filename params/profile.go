package params

import "mymodule/entity"

type ProfileRequest struct {
	Id uint `json:"id"`
}

type ProfileResponse struct {
	User entity.User `json:"user"`
}
