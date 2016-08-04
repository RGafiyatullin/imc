package main

import (
	"fmt"
	"github.com/rgafiyatullin/imc/server/actor"
	"github.com/rgafiyatullin/imc/server/netsrv"
)

func main() {
	fmt.Println("Helloes! I'm the IMC daemon")

	topActorCtx := actor.NewCtx()
	topActorCtx.Log().Info("System start")

	listener, _ := netsrv.StartListener(topActorCtx.NewChild("listener"), ":6379")

	listener.Join()
}
