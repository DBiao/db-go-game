package main

import (
	"db-go-game/pkg/commands"
	"db-go-game/services/admin/internal/server"
)

func main() {
	commands.Run(server.NewServer())
}
