package controllers

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/endpoint"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("controllers")
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	endpoint endpoint.IEndpoint
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	endpoint endpoint.IEndpoint) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		endpoint: endpoint,
	}
}
