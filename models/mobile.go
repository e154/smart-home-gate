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

package models

import (
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"github.com/e154/smart-home-gate/system/uuid"
	"time"
)

type Mobile struct {
	Id        string    `json:"id"`
	Token     uuid.UUID `json:"token"`
	ServerId  string    `json:"server_id"`
	RequestId string    `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Mobile) GenAccessToken() (token string) {

	var timestamp = s.CreatedAt.Unix()
	token = fmt.Sprintf("%s-%s-%d-%s", s.Id, s.RequestId, timestamp, common.Sha256(s.RequestId+s.Token.String()+fmt.Sprintf("%d", timestamp)))

	return
}
