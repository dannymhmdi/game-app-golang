package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mymodule/contract/golang/presence"
)

func main() {
	conn, dErr := grpc.Dial(":8086", grpc.WithInsecure())
	if dErr != nil {
		panic(dErr)
	}

	defer conn.Close()

	presenceClient := presence.NewPresenceServiceClient(conn)
	res, gErr := presenceClient.GetPresence(context.Background(), &presence.GetPresenceRequest{
		UserIds: []uint64{1, 2, 3, 4},
	})

	if gErr != nil {
		panic(gErr)
	}

	fmt.Printf("grpc request response:%+v\n", res.OnlinePlayers)
}
