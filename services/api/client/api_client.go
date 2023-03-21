package client

import (
	"context"
	"db-go-game/pkg/common/dgrpc"
	"db-go-game/pkg/conf"
	"db-go-game/pkg/proto/logic"
	"google.golang.org/grpc"
)

type IApiClient interface {
	KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error)
}

type apiClient struct {
	opt *dgrpc.ClientDialOption
}

func NewAuthClient(etcd *conf.Etcd, server *conf.GrpcServer, jaeger *conf.Jaeger, clientName string) IApiClient {
	return &apiClient{dgrpc.NewClientDialOption(etcd, server, jaeger, clientName)}
}

func (a *apiClient) GetClientConn() (conn *grpc.ClientConn) {
	conn = dgrpc.GetClientConn(a.opt.DialOption)
	return
}

func (a *apiClient) KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error) {

	return &logic.KickOffLineResp{Code: 1000}, nil
}
