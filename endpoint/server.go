// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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
