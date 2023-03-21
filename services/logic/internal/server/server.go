package server

import (
	"db-go-game/pkg/commands"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/common/dmysql"
	"db-go-game/pkg/common/dredis"
	"db-go-game/services/logic/internal/config"
	"flag"
)

var portConf = flag.Int("p", 18080, "api default listen port 18080")

type server struct {
	grpcServer IGrpcServer
}

func NewServer(grpcServer IGrpcServer) commands.MainInstance {
	return &server{grpcServer: grpcServer}
}

func (s *server) Initialize() (err error) {
	conf := config.GetConfig()
	dlog.Shared(conf.Log, conf.Name)
	dmysql.NewMysqlClient(conf.Mysql)
	dredis.NewRedisClient(conf.Redis)
	return
}

func (s *server) RunLoop() {
	s.grpcServer.Run()
}

func (s *server) Destroy() {

}
