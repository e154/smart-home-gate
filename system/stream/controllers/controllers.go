package controllers

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("stream.controllers")
)

type StreamControllers struct {
	Client *ControllerClient
}

func NewStreamControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService) *StreamControllers {

	ctrls := &StreamControllers{
		Client: NewControllerClient(adaptors),
	}

	// client
	stream.Subscribe("get_token", ctrls.Client.GetToken)

	return ctrls
}
