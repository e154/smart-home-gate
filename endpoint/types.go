package endpoint

import (
	m "github.com/e154/smart-home-gate/models"
)

type IEndpoint interface {
	IEndpointMobile
	IEndpointServer
}

type IEndpointMobile interface {
	RegisterMobile(server *m.Server) (string, error)
	RemoveMobileToken(*m.Server, string) error
	ListMobileToken(limit, offset int64) (list []*m.Mobile, total int64, err error)
}

type IEndpointServer interface {
	RegisterServer() (string, error)
	RemoveServerToken(*m.Server) error
}
