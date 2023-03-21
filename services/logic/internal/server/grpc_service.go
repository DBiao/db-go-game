package server

import (
	"db-go-game/pkg/common/dgrpc"
	"db-go-game/pkg/proto/logic"
	"db-go-game/services/logic/internal/config"
	"db-go-game/services/logic/internal/service"
	"google.golang.org/grpc"
	"io"
)

type IGrpcServer interface {
	Run()
}

type grpcServer struct {
	logic.UnimplementedApiServer
	cfg         *config.Config
	apiService  service.IApiService
	lGrpcServer *dgrpc.GrpcServer
}

func NewGrpcServer(apiService service.IApiService) IGrpcServer {
	return &grpcServer{apiService: apiService}
}

func (s *grpcServer) Run() {
	var (
		srv    *grpc.Server
		closer io.Closer
	)
	srv, closer = dgrpc.NewServer(s.cfg.GrpcServer)
	defer func() {
		if closer != nil {
			closer.Close()
		}
	}()

	logic.RegisterApiServer(srv, s)
	s.lGrpcServer = dgrpc.NewGrpcServer(s.cfg.GrpcServer, s.cfg.Etcd)
	s.lGrpcServer.RunServer(srv)
}
