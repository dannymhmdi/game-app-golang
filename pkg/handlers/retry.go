package handlers

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"mymodule/contract/broker"
	"mymodule/contract/golang/matchingPlayer"
	"mymodule/entity"
	"time"
)

func Retry(attempts int, initialDelay time.Duration, fn broker.PublisherFunc, rollBackFn func(userId uint, category entity.Category) error, event string, payload string) {
	delay := initialDelay
	var err error
	for i := 0; i < attempts; i++ {
		err = fn(event, payload)
		if err == nil {

			return
		}

		if i < attempts-1 {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}

	//after specified attempt to publish msg roll back users to waitinglist again

	decodedPayload, dErr := base64.StdEncoding.DecodeString(payload)
	if dErr != nil {
		fmt.Println("Error decoding payload:(handlers.Retry)", dErr)

		return
	}
	var matchedUsers matchingPlayer.MatchedPlayers
	if uErr := proto.Unmarshal(decodedPayload, &matchedUsers); uErr != nil {
		fmt.Println("Error decoding payload:(handlers.Retry)", uErr)
	}

	for _, v := range matchedUsers.UserIds {
		if rErr := rollBackFn(uint(v), entity.Category(matchedUsers.Category)); rErr != nil {
			fmt.Println("rErr(handlers.Retry):", rErr)

			return
		}
	}

}
