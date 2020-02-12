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
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/api/server/controllers"
	"github.com/e154/smart-home-gate/api/websocket"
	"github.com/e154/smart-home-gate/endpoint"
	"github.com/e154/smart-home-gate/system/config"
	"github.com/e154/smart-home-gate/system/dig"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/e154/smart-home-gate/system/logging"
	"github.com/e154/smart-home-gate/system/metrics"
	"github.com/e154/smart-home-gate/system/migrations"
	"github.com/e154/smart-home-gate/system/orm"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/e154/smart-home-gate/system/stream_proxy"
)

func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(server.NewServer)
	container.Provide(server.NewServerConfig)
	container.Provide(controllers.NewControllers)
	container.Provide(config.ReadConfig)
	container.Provide(graceful_service.NewGracefulService)
	container.Provide(graceful_service.NewGracefulServicePool)
	container.Provide(graceful_service.NewGracefulServiceConfig)
	container.Provide(orm.NewOrm)
	container.Provide(orm.NewOrmConfig)
	container.Provide(migrations.NewMigrations)
	container.Provide(migrations.NewMigrationsConfig)
	container.Provide(adaptors.NewAdaptors)
	container.Provide(logging.NewLogrus)
	container.Provide(logging.NewLogBackend)
	container.Provide(stream.NewStreamService)
	container.Provide(stream.NewHub)
	container.Provide(endpoint.NewEndpoint)
	container.Provide(websocket.NewWebSocket)
	container.Provide(stream_proxy.NewStreamProxy)
	container.Provide(metrics.NewMetricConfig)
	container.Provide(metrics.NewMetricServer)

	return
}
