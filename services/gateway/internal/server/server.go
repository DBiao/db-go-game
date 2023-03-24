package server

import (
	"db-go-game/pkg/commands"
)

type server struct {
}

func NewServer() commands.MainInstance {
	return &server{}
}

func (s *server) Initialize() (err error) {
	return
}

func (s *server) RunLoop() {
}

func (s *server) Destroy() {
}
