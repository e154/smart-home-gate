package endpoint

import (
	"github.com/e154/smart-home-gate/adaptors"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("endpoint")
)

type Endpoint struct {
	*Mobile
	*Server
}

func NewEndpoint(adaptors *adaptors.Adaptors) IEndpoint {
	common := NewCommonEndpoint(adaptors)
	return &Endpoint{
		Mobile: NewMobile(common),
		Server: NewServer(common),
	}
}
