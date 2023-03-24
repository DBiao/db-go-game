package client

import (
	"context"
	"db-go-game/pkg/common/dgrpc"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/conf"
	"db-go-game/pkg/proto/logic"
	"errors"
	"google.golang.org/grpc"
)

type IApiClient interface {
	KickOffLine(req *logic.KickOffLineReq) (*logic.KickOffLineResp, error)
}

type apiClient struct {
	opt *dgrpc.ClientDialOption
}

func NewApiClient(etcd *conf.Etcd, server *conf.GrpcServer, jaeger *conf.Jaeger, clientName string) IApiClient {
	return &apiClient{dgrpc.NewClientDialOption(etcd, server, jaeger, clientName)}
}

func (a *apiClient) GetClientConn() (conn *grpc.ClientConn) {
	conn = dgrpc.GetClientConn(a.opt.DialOption)
	return
}

func (a *apiClient) KickOffLine(req *logic.KickOffLineReq) (*logic.KickOffLineResp, error) {
	conn := a.GetClientConn()
	if conn == nil {
		return &logic.KickOffLineResp{Code: 1000}, errors.New("")
	}
	client := logic.NewApiClient(conn)
	var err error
	resp, err := client.KickOffLine(context.Background(), req)
	if err != nil {
		dlog.Error(err.Error())
		return &logic.KickOffLineResp{Code: 1000}, err
	}
	return resp, nil
}
