package main

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/api/server"
	"github.com/e154/smart-home-gate/api/server/controllers"
	"github.com/e154/smart-home-gate/endpoint"
	"github.com/e154/smart-home-gate/system/config"
	"github.com/e154/smart-home-gate/system/dig"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/e154/smart-home-gate/system/logging"
	"github.com/e154/smart-home-gate/system/migrations"
	"github.com/e154/smart-home-gate/system/orm"
	"github.com/e154/smart-home-gate/system/stream"
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

	return
}
