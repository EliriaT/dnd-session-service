package server

import (
	"github.com/EliriaT/dnd-user-service/config"
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  config.Config
}

func NewServer(queries *db.Queries, config config.Config) (*Server, error) {
	server := &Server{
		queries: queries,
		config:  config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/sessions", server.getCampaignSessions)
	router.POST("/sessions/create", server.createSession)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
