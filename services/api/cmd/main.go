package main

import (
	"db-go-game/pkg/commands"
	"db-go-game/services/api/internal/server"
)

func main() {
	commands.Run(server.NewServer())
}
