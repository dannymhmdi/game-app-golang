package dto

import "mymodule/entity"

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User `json:"user"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Token   Token  `json:"token"`
}

type ProfileRequest struct {
	Id uint `json:"id"`
}

type ProfileResponse struct {
	User entity.User `json:"user"`
}
