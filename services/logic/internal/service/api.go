package service

import (
	"context"
	"db-go-game/pkg/proto/logic"
	"google.golang.org/grpc"
)

var count int

type IApiService interface {
	KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error)
}

type apiService struct {
}

func NewAuthService() IApiService {
	return &apiService{}
}

func (a *apiService) KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error) {
	count++
	return &logic.KickOffLineResp{Code: int32(count)}, nil
}
