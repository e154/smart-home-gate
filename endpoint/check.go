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
	"fmt"
	"github.com/e154/smart-home-gate/common"
	"strconv"
	"strings"
)

type Check struct {
	*CommonEndpoint
}

func NewCheck(common *CommonEndpoint) *Check {
	return &Check{CommonEndpoint: common}
}

func (c *Check) CheckMobileAccessToken(accessToken string) (ok bool) {

	data := strings.Split(accessToken, "-")
	if len(data) != 4 {
		return
	}

	mobileClientId := data[0]
	requestRandomId := data[1]
	hash := data[3]

	timestamp, errw := strconv.Atoi(data[2])
	if errw != nil {
		return
	}

	if len(requestRandomId) < 8 {
		return
	}

	mobileObj, err := c.adaptors.Mobile.GetById(mobileClientId)
	if err != nil {
		return
	}

	if hash != common.Sha256(requestRandomId+mobileObj.Token.String()+fmt.Sprintf("%d", timestamp)) {
		return
	}

	ok = true

	return
}
