package websocket

import (
	"github.com/e154/smart-home-gate/adaptors"
	. "github.com/e154/smart-home-gate/api/websocket/controllers"
	"github.com/e154/smart-home-gate/endpoint"
	"github.com/e154/smart-home-gate/system/graceful_service"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("websocket")
)


type WebSocket struct {
	Controllers *Controllers
}

func NewWebSocket(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint endpoint.IEndpoint,
	graceful *graceful_service.GracefulService) *WebSocket {

	server := &WebSocket{
		Controllers: NewControllers(adaptors, stream, endpoint),
	}

	graceful.Subscribe(server)

	return server
}

func (s *WebSocket) Start() {
	log.Infof("Serving websocket service")
}

func (s *WebSocket) Shutdown() {
	log.Info("Server exiting")
}