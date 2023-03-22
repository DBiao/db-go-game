package client

import (
	"context"
	"db-go-game/pkg/common/dgrpc"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/conf"
	"db-go-game/pkg/proto/logic"
	"google.golang.org/grpc"
)

type ILogicClient interface {
	KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error)
}

type logicClient struct {
	opt *dgrpc.ClientDialOption
}

func NewLogicClient(etcd *conf.Etcd, server *conf.GrpcServer, jaeger *conf.Jaeger, clientName string) ILogicClient {
	return &logicClient{dgrpc.NewClientDialOption(etcd, server, jaeger, clientName)}
}

func (a *logicClient) GetClientConn() (conn *grpc.ClientConn) {
	conn = dgrpc.GetClientConn(a.opt.DialOption)
	return
}

func (a *logicClient) KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error) {
	conn := a.GetClientConn()
	if conn == nil {
		return &logic.KickOffLineResp{Code: 1000}, nil
	}
	client := logic.NewApiClient(conn)
	var err error
	resp, err := client.KickOffLine(context.Background(), in)
	if err != nil {
		dlog.Warn(err.Error())
	}
	return resp, nil
}
