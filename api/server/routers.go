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
