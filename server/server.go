package server

import (
	"github.com/EliriaT/dnd-user-service/config"
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/gin-gonic/gin"
	//"golang.org/x/net/websocket"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  config.Config
	//sessionConns map[int64]map[*websocket.Conn]bool
}

func NewServer(queries *db.Queries, config config.Config) (*Server, error) {
	//sockets := make(map[int64]map[*websocket.Conn]bool)
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
	router.POST("/sessions/objects", server.editObject)
	router.POST("/sessions/characters", server.editCharacter)
	router.GET("/sessions/:sessionId/map-state", server.getSessionMapState)
	router.GET("/sessions/:sessionId", server.getSessionById)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
