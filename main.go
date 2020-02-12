// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/api/websocket"
	"github.com/e154/smart-home-gate/system/graceful_service"
	l "github.com/e154/smart-home-gate/system/logging"
	"github.com/e154/smart-home-gate/system/metrics"
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
	err := container.Invoke(func(server *server.Server,
		graceful *graceful_service.GracefulService,
		back *l.LogBackend,
		ws *websocket.WebSocket,
		streamProxy *stream_proxy.StreamProxy,
		metric *metrics.MetricServer) {

		l.Initialize(back)
		go server.Start()
		go ws.Start()
		go streamProxy.Start()
		go metric.Start()

		graceful.Wait()
	})

	if err != nil {
		panic(err.Error())
	}
}
