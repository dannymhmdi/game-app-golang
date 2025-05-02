package presenceserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mymodule/adaptor/redis"
	"mymodule/config"
	"mymodule/contract/golang/presence"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/pkg/slice"
	"mymodule/repository/redis/redisPresence"
	"mymodule/service/presenceService"
	"net"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	presenceSvc presenceService.Service
}

func New(presenceSvc presenceService.Service) *Server {
	return &Server{
		//UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		presenceSvc: presenceSvc,
	}
}

func (s Server) Start() {
	address := fmt.Sprintf("localhost:%d", 8086)
	listener, lErr := net.Listen("tcp", address)
	if lErr != nil {
		log.Fatalf("failed to listen: %v", lErr)
	}

	appConfig := config.Load()

	redisAdaptor := redis.New(appConfig.RedisConfig)
	pRepo := redisPresence.New(redisAdaptor, appConfig.RedisPresence)
	p := presenceService.New(pRepo)
	presenceSvc := New(*p)

	grpcServer := grpc.NewServer()
	presence.RegisterPresenceServiceServer(grpcServer, presenceSvc)
	fmt.Printf("presence grpc server start on %s\n", address)
	if sErr := grpcServer.Serve(listener); sErr != nil {
		log.Fatalf("failed to serve: %v", sErr)
	}

}

//call GetPresence function from client

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
