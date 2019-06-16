package endpoint

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("endpoint")
)

type Endpoint struct {
}

func NewEndpoint(adaptors *adaptors.Adaptors) *Endpoint {
	//common := NewCommonEndpoint(adaptors)
	return &Endpoint{

	}
}
