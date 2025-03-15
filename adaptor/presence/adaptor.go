package presence

import (
	"context"
	"google.golang.org/grpc"
	"mymodule/contract/golang/presence"
	"mymodule/params"
	"mymodule/pkg/slice"
)

type PresenceClient struct {
	client presence.PresenceServiceClient
}

func New(client *grpc.ClientConn) *PresenceClient {
	return &PresenceClient{
		client: presence.NewPresenceServiceClient(client),
	}

}

//we want to pass PresenceClient adaptor as getPresenceClient field in matchMakingSVc so we should implement its method

func (p PresenceClient) GetPresence(ctx context.Context, req params.GetPresenceRequest) (params.GetPresenceResponse, error) {
	//call GetPresence method which writes in remote presence-grpc-server
	res, gErr := p.client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.UintToUint64Mapper(req.UserIDs)})
	if gErr != nil {
		return params.GetPresenceResponse{}, gErr
	}

	return params.GetPresenceResponse{OnlinePlayers: slice.OnlinePlayerMapperToParams(res.OnlinePlayers)}, nil

}
