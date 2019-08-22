package endpoint

import (
	"github.com/e154/smart-home-gate/adaptors"
)

type CommonEndpoint struct {
	adaptors *adaptors.Adaptors
}

func NewCommonEndpoint(adaptors *adaptors.Adaptors) *CommonEndpoint {
	return &CommonEndpoint{
		adaptors: adaptors,
	}
}