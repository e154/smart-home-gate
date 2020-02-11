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

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerCheck struct {
	*ControllerCommon
}

func NewControllerCheck(common *ControllerCommon) *ControllerCheck {
	return &ControllerCheck{ControllerCommon: common}
}

// swagger:operation GET /check/mobile_access_token check
// ---
// summary: mobile access token check page
// description:
// security:
// - ServerAuthorization: []
// consumes:
// - text/plain
// produces:
// - text/plain
// tags:
// - check
// responses:
//   "200":
//	   description: Success response
//   "400":
//	   description: Bad request
//   "404":
//	   description: Not found
//
func (i ControllerCheck) CheckMobileAccessToken(ctx *gin.Context) {

	accessToken := ctx.GetHeader("ServerAuthorization")
	if accessToken == "" {
		ctx.String(http.StatusBadRequest, "need access token")
		return
	}

	ok := i.endpoint.CheckMobileAccessToken(accessToken)
	fmt.Println("ok", ok)
	if ok {
		ctx.String(http.StatusOK, "ok")
	} else {
		ctx.String(http.StatusNotFound, "not found")
	}

	return
}

// swagger:operation GET /check/mobile_access check
// ---
// summary: mobile access connection page
// description:
// consumes:
// - text/plain
// produces:
// - text/plain
// tags:
// - check
// responses:
//   "200":
//	   description: Success response
//
func (i ControllerCheck) CheckMobileAccess(ctx *gin.Context) {
	apiVersion := "smart-home-gate"
	ctx.String(http.StatusOK, apiVersion)
	return
}
