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
	"errors"
	"github.com/e154/smart-home-gate/common"
	m "github.com/e154/smart-home-gate/models"
	"github.com/e154/smart-home-gate/system/uuid"
)

type Mobile struct {
	*CommonEndpoint
}

func NewMobile(common *CommonEndpoint) *Mobile {
	return &Mobile{CommonEndpoint: common}
}

func (c *Mobile) RegisterMobile(server *m.Server) (token string, err error) {

	mobileClient := &m.Mobile{
		Token:     uuid.NewV4(),
		ServerId:  server.Id,
		RequestId: common.RandomString(10, common.Charset),
	}

	if _, err = c.adaptors.Mobile.Add(mobileClient); err != nil {
		return
	}

	token = mobileClient.GenAccessToken()

	return
}

func (c *Mobile) RemoveMobileToken(server *m.Server, token string) (err error) {

	mobile, err := c.adaptors.Mobile.GetByAccessToken(token)
	if err != nil {
		return
	}

	var exist bool
	for _, m := range server.Mobiles {
		if m.Token == mobile.Token {
			exist = true
		}
	}

	if !exist {
		err = errors.New("mobile not found")
		return
	}

	err = c.adaptors.Mobile.Remove(mobile)

	return
}

func (c *Mobile) ListMobileToken(limit, offset int64) (list []*m.Mobile, total int64, err error) {
	list, total, err = c.adaptors.Mobile.List(limit, offset)
	return
}
