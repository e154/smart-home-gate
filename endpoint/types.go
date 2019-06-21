package endpoint

import (
	m "github.com/e154/smart-home-gate/models"
)

type IEndpoint interface {
	IEndpointMobile
	IEndpointServer
}

type IEndpointMobile interface {
	RegisterMobile() (string, error)
	RemoveMobileToken(*m.Mobile) error
}

type IEndpointServer interface {
	RegisterServer() (string, error)
	RemoveServerToken(*m.Server) error
}
