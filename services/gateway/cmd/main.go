package main

import (
	"db-go-game/pkg/commands"
	"db-go-game/services/gateway/dig"
)

func main() {
	dig.Invoke(func(svc commands.MainInstance) {
		commands.Run(svc)
	})
}
