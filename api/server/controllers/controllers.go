package controllers

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/e154/smart-home-gate/endpoint"
)

type Controllers struct {
	Index   *ControllerIndex
	Swagger *ControllerSwagger
	Check   *ControllerCheck
}

func NewControllers(adaptors *adaptors.Adaptors,
	endpoint endpoint.IEndpoint) *Controllers {
	common := NewControllerCommon(adaptors, endpoint)
	return &Controllers{
		Index:   NewControllerIndex(common),
		Swagger: NewControllerSwagger(common),
		Check:   NewControllerCheck(common),
	}
}
