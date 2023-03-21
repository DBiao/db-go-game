package service

import (
	"context"
	"db-go-game/pkg/proto/logic"
	"db-go-game/services/logic/internal/config"
	"google.golang.org/grpc"
)

type IApiService interface {
	KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error)
}

type apiService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) IApiService {
	return &apiService{cfg: cfg}
}

func (a *apiService) KickOffLine(ctx context.Context, in *logic.KickOffLineReq, opts ...grpc.CallOption) (*logic.KickOffLineResp, error) {

	return &logic.KickOffLineResp{Code: 1000}, nil
}
