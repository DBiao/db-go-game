package server

import (
	"db-go-game/pkg/commands"
	"db-go-game/pkg/common/dgin"
	"db-go-game/pkg/common/dlog"
	"db-go-game/services/admin/internal/config"
	"db-go-game/services/admin/internal/router"
	"flag"
)

var portConf = flag.Int("p", 18080, "api default listen port 18080")

type server struct {
	ginServer *dgin.GinServer
}

func NewServer() commands.MainInstance {
	return &server{}
}

func (s *server) Initialize() (err error) {
	conf := config.GetConfig()
	dlog.Shared(conf.Log, conf.Name)
	s.ginServer = dgin.NewGinServer()
	router.Register(s.ginServer.Engine)
	return
}

func (s *server) RunLoop() {
	conf := config.GetConfig()
	s.ginServer.Run(conf.Port)
}

func (s *server) Destroy() {

}
