package notificationservice

import (
	"context"
	"errors"
	"fmt"
	"mymodule/contract/broker"
	"mymodule/contract/golang/matchingPlayer"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/slice"
	"time"

	"github.com/kavenegar/kavenegar-go"
)

type Service struct {
	msgConsumer broker.Consumer
} 



func (s Service) SendNotifToUser (ctx context.Context,req params.NotficationRequest) (params.NotficationResponse,error) {
	done :=make(chan bool)
	deliveryMsg:= make(chan matchingPlayer.MatchedPlayers)
rabbitConn,rabbitchan:=s.msgConsumer.Consume(ctx,"notification_queue",done,deliveryMsg)

go func ()  {
	for msg:= range deliveryMsg {
		game := entity.Game{
			ID:        0,
			Category:  entity.Category(msg.Category),
			PlayersID: slice.Uint64ToUintMapper(msg.UserIds),
			StartTime: time.Now(),
		}

		fmt.Println("game",game)
		
		api := kavenegar.New("4134725861443144367A3235584F6F647077556E36363365505A715355516B3568434B37335A32674F65343D")
		sender := "2000660110"
		receptor := []string{"09127275236"}
		message := "وب سرویس پیام کوتاه کاوه نگار"
		if _, err := api.Message.Send(sender, receptor, message, nil); err != nil {
			var apiErr *kavenegar.APIError
			var httpErr *kavenegar.HTTPError
			if errors.As(err, &apiErr) {
				fmt.Println("api error occurred", err.Error())
				done<-false
				continue
			} else if errors.As(err, &httpErr) {
				done<-false
				fmt.Println("http error occurred", err.Error())
				continue
			} else {
				done<-false
				fmt.Println("error", err.Error())
				continue
			}
	
		} 
	
	
		done <- true
	}
}()

return params.NotficationResponse{
RabbitConnection: rabbitConn,
RabbitChannel: rabbitchan,
},nil
}