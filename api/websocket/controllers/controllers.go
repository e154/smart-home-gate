package controllers

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/endpoint"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("stream.controllers")
)

type Controllers struct {
	Server *ControllerServer
	Mobile *ControllerMobile
}

func NewControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint endpoint.IEndpoint) *Controllers {
	common := NewControllerCommon(adaptors, stream, endpoint)
	return &Controllers{
		Server: NewControllerServer(common, stream),
		Mobile: NewControllerMobile(common, stream),
	}
}
