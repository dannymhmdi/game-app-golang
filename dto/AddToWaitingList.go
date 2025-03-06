package dto

import (
	"mymodule/entity"
	"time"
)

type AddToWaitingListRequest struct {
	UserId   uint            `json:"user_id"`
	Category entity.Category `json:"category"`
}

type AddToWaitingListResponse struct {
	Message string `json:"message"`
	//Timeout time.Duration `json:"timeout"`
	Timeout time.Time `json:"timeout"`
}
