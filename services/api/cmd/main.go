package main

import (
	"db-go-game/pkg/commands"
	"db-go-game/services/api/dig"
)

func main() {
	dig.Invoke(func(srv commands.MainInstance) {
		commands.Run(srv)
	})
}
