package main

import (
	"errors"
	"fmt"
	"github.com/kavenegar/kavenegar-go"
)

func main() {
	api := kavenegar.New("4134725861443144367A3235584F6F647077556E36363365505A715355516B3568434B37335A32674F65343D")
	sender := "2000660110"
	receptor := []string{"09127275236"}
	message := "وب سرویس پیام کوتاه کاوه نگار"
	if res, err := api.Message.Send(sender, receptor, message, nil); err != nil {
		var apiErr *kavenegar.APIError
		var httpErr *kavenegar.HTTPError
		if errors.As(err, &apiErr) {
			fmt.Println("api error occurred", err.Error())
		} else if errors.As(err, &httpErr) {
			fmt.Println("http error occurred", err.Error())
		} else {
			fmt.Println("error", err.Error())
		}

	} else {
		for _, r := range res {
			fmt.Println("MessageID   = ", r.MessageID)
			fmt.Println("Status      = ", r.Status)
			fmt.Println("message", r.Message)
			//...
		}
	}
}
