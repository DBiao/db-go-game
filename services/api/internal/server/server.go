package server

import (
	"db-go-game/pkg/commands"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/common/dmysql"
	"db-go-game/pkg/common/dredis"
	"db-go-game/pkg/common/xgin"
	"db-go-game/services/api/internal/config"
	"db-go-game/services/api/internal/router"
	"flag"
)

var portConf = flag.Int("p", 18080, "gateway default listen port 10080")

type server struct {
	ginServer *xgin.GinServer
}

func NewServer(gin *xgin.GinServer) commands.MainInstance {
	return &server{ginServer: gin}
}

func (s *server) Initialize() (err error) {
	conf := config.GetConfig()
	dlog.Shared(conf.Log, conf.Name)
	dmysql.NewMysqlClient(conf.Mysql)
	dredis.NewRedisClient(conf.Redis)
	router.Register(s.ginServer.Engine)
	return
}

func (s *server) RunLoop() {
	conf := config.GetConfig()
	s.ginServer.Run(conf.Port)
}

func (s *server) Destroy() {

}
