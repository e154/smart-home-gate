// Smart home gate
//
// This documentation describes APIs found under https://github.com/e154/smart-home-gate
//
//     BasePath: /
//     Version: 1.0.0
//     License: MIT https://raw.githubusercontent.com/e154/smart-home-gate/master/LICENSE
//     Contact: Alex Filippov <support@e154.ru> https://e154.github.io/smart-home-gate/
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - ApiKeyAuth
//
//     SecurityDefinitions:
//     ApiKeyAuth:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package server

import (
	"context"
	"fmt"
	"github.com/e154/smart-home-gate/api/server/controllers"
	"github.com/e154/smart-home-gate/system/stream"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"time"
)

var (
	log = logging.MustGetLogger("server")
)

type Server struct {
	Config        *ServerConfig
	Controllers   *controllers.Controllers
	engine        *gin.Engine
	server        *http.Server
	logger        *ServerLogger
	streamService *stream.StreamService
}

func (s *Server) Start() {

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port),
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Serving server at http://[::]:%d", s.Config.Port)
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
	log.Info("Server exiting")
}

func NewServer(cfg *ServerConfig,
	controllers *controllers.Controllers,
	streamService *stream.StreamService) (newServer *Server) {

	logger := &ServerLogger{log}

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	newServer = &Server{
		Config:        cfg,
		Controllers:   controllers,
		engine:        engine,
		logger:        logger,
		streamService: streamService,
	}

	newServer.setControllers()

	return
}
