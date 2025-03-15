# GRPC

If we need to implement microservice we use grpc call to comunicate between services
we implement grpc server in delivery like this project and grpc client in adaptor.
we pass this adoptor to matchmaking service as GetPresenceClient interface because this adoptor implement all
this interface method so can use as this interface .


### guide

we use the generated go code by protobuf in presence_grpc.pb.go file for example:

in PresenceServiceServer interface we understood that we should implement all PresenceServiceServer interface
so we implement the methods for grpcServer and we inject presence service directly to grpcServer to use its methods like 
```go

type Server struct {
presence.UnimplementedPresenceServiceServer
presenceSvc presenceService.Service
}

func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	res, gErr := s.presenceSvc.GetPresence(ctx, params.GetPresenceRequest{UserIDs: slice.Uint64ToUintMapper(req.UserIds)})
	fmt.Printf("grpc req:%+v\n", req)
	if gErr != nil {
		return nil, richerr.New().
			SetMsg(gErr.Error()).
			SetOperation("presenceserver.GetPresence").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(gErr)

	}

	return &presence.GetPresenceResponse{OnlinePlayers: slice.OnlinePlayerMapperToProtobuf(res.OnlinePlayers)}, nil
}
```
 after that implement adaptor as grpc-client to request to grpc server to call
 GetPresence method of grpc server. like:
 ```go
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
 ```

