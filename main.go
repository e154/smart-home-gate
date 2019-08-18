package main

import (
	"fmt"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/api/websocket"
	"github.com/e154/smart-home-gate/system/graceful_service"
	l "github.com/e154/smart-home-gate/system/logging"
	"github.com/e154/smart-home-gate/system/migrations"
	"github.com/e154/smart-home-gate/system/stream_proxy"
	"github.com/op/go-logging"
	"os"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(shortVersionBanner, GetHumanVersion())
			return
		default:
			fmt.Printf(verboseVersionBanner, "v1", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(shortVersionBanner, "")

	container := BuildContainer()
	err := container.Invoke(func(m *migrations.Migrations) {
		m.Up()
	})

	if err != nil {
		panic(err.Error())
	}

	err = container.Invoke(func(server *server.Server,
		graceful *graceful_service.GracefulService,
		back *l.LogBackend,
		ws *websocket.WebSocket,
		streamProxy *stream_proxy.StreamProxy) {

		l.Initialize(back)
		go server.Start()
		go ws.Start()
		go streamProxy.Start()

		graceful.Wait()
	})

	if err != nil {
		panic(err.Error())
	}
}
