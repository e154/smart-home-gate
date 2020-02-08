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

package server

import (
	"github.com/e154/smart-home-gate/system/swaggo/gin-swagger/swaggerFiles"
	"github.com/gin-gonic/gin"
)

func (s *Server) setControllers() {

	r := s.engine

	basePath := r.Group("/")

	basePath.GET("/", s.Controllers.Index.Index)
	basePath.GET("/swagger", func(context *gin.Context) {
		context.Redirect(302, "/swagger/index.html")
	})
	basePath.GET("/swagger/*any", s.Controllers.Swagger.WrapHandler(swaggerFiles.Handler))

	// check
	basePath.GET("/check/mobile_access_token", s.Controllers.Check.CheckMobileAccessToken)
	basePath.GET("/check/mobile_access", s.Controllers.Check.CheckMobileAccess)

	// ws
	basePath.GET("/ws", s.streamService.Ws)
	basePath.GET("/ws/*any", s.streamService.Ws)

}
