package endpoint

import (
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/uuid"
)

type Server struct {
	*CommonEndpoint
}

func NewServer(common *CommonEndpoint) *Server {
	return &Server{CommonEndpoint: common}
}

func (s *Server) RegisterServer() (server *m.Server, err error) {

	server = &m.Server{
		Token: uuid.NewV4().String(),
	}

	if server.Id, err = s.adaptors.Server.Add(server); err != nil {
		return
	}

	return
}

func (s *Server) RemoveServerToken(server *m.Server) (err error) {
	err = s.adaptors.Server.Remove(server)
	return
}
