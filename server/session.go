package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/EliriaT/dnd-user-service/server/dto"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (server *Server) getCampaignSessions(ctx *gin.Context) {
	var req dto.GetSessionsByCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessions, err := server.queries.GetSessionsByCampaignAndCharacter(ctx, db.GetSessionsByCampaignAndCharacterParams{
		CampaignID:  req.CampaignId,
		CharacterID: req.UserId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sessions == nil {
		sessions = []db.GetSessionsByCampaignAndCharacterRow{}
	}

	ctx.JSON(http.StatusOK, sessions)
}

func (server *Server) createSession(ctx *gin.Context) {
	var req dto.CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	session, err := server.queries.CreateSession(ctx, db.CreateSessionParams{
		Name:       req.Name,
		CampaignID: req.CampaignID,
		MapID:      req.MapID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = server.queries.AddSessionAllowedCharacter(ctx, db.AddSessionAllowedCharacterParams{
		SessionID: session.ID,
		Column2:   req.AllowedChars,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": session.ID})
}

func (server *Server) editCharacter(ctx *gin.Context) {
	var req dto.EditCharacterPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.queries.UpsertCharacterPosition(ctx, db.UpsertCharacterPositionParams{
		SessionID:        req.SessionID,
		CharID:           req.CharacterID,
		XPos:             int32(req.X),
		YPos:             int32(req.Y),
		IsVisible:        true,
		ModificationDate: time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (server *Server) editObject(ctx *gin.Context) {
	var req dto.EditObjectPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.queries.UpsertObjectPosition(ctx, db.UpsertObjectPositionParams{
		SessionID:        req.SessionID,
		ObjectID:         req.ObjectID,
		XPos:             int32(req.X),
		YPos:             int32(req.Y),
		IsVisible:        *req.IsVisible,
		ModificationDate: time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (server *Server) getSessionMapState(ctx *gin.Context) {
	var uri dto.GetIdRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing id parameter"})
		return
	}

	chars, err := server.queries.GetCharactersBySession(ctx, uri.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get characters: " + err.Error()})
		return
	}

	objs, err := server.queries.GetObjectsBySession(ctx, uri.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get objects: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"characters": chars,
		"objects":    objs,
	})
}

func (server *Server) getSessionById(ctx *gin.Context) {
	var uri dto.GetIdRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing id parameter"})
		return
	}

	session, err := server.queries.GetSessionByID(ctx, uri.SessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get session: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

func (server *Server) connectToSession(ctx *gin.Context) {
	var uri dto.GetIdRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing id parameter"})
		return
	}
	isDMQuery := ctx.Query("isDM")
	if isDMQuery != "" {
		isDM, err := strconv.ParseBool(isDMQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing isDM parameter"})
			return
		}

		if isDM {
			err := server.queries.SetSessionActive(ctx, uri.SessionID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	handler := func(ws *websocket.Conn) {
		server.WebSocketHandler(ws, uri.SessionID)
	}

	websocket.Handler(handler).ServeHTTP(ctx.Writer, ctx.Request)
}

func (server *Server) WebSocketHandler(ws *websocket.Conn, sessionId int64) {
	defer ws.Close()

	log.Println("new connection established for session:", sessionId)
	log.Println("new incoming connection from client", ws.RemoteAddr())

	server.mutex.Lock()
	if server.sessionConns[sessionId] == nil {
		server.sessionConns[sessionId] = make(map[*websocket.Conn]bool)
	}
	server.sessionConns[sessionId][ws] = true
	server.mutex.Unlock()

	chars, err := server.queries.GetCharactersBySession(context.Background(), sessionId)
	if err != nil {
		log.Println(err)
		ws.Write([]byte("failed to load characters"))
		ws.Close()
		return
	}

	objs, err := server.queries.GetObjectsBySession(context.Background(), sessionId)
	if err != nil {
		ws.Write([]byte("failed to load objects"))
		ws.Close()
		return
	}

	msg, _ := json.Marshal(gin.H{"characters": chars, "objects": objs})
	ws.Write(msg)

	server.readLoop(ws, sessionId)
}

func (server *Server) readLoop(ws *websocket.Conn, sessionId int64) {
	for {
		var raw string
		if err := websocket.Message.Receive(ws, &raw); err != nil {
			if err == io.EOF {
				delete(server.sessionConns[sessionId], ws)
				break
			}
			log.Println("Received error:", err)
			break
		}

		var msg dto.WebsocketMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		switch msg.Type {
		case "editCharacter":
			var payload dto.EditCharacterPositionRequest
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				log.Println("Invalid payload for editCharacter:", err)
				continue
			}
			err := server.queries.UpsertCharacterPosition(context.Background(), db.UpsertCharacterPositionParams{
				SessionID:        payload.SessionID,
				CharID:           payload.CharacterID,
				XPos:             int32(payload.X),
				YPos:             int32(payload.Y),
				IsVisible:        true,
				ModificationDate: time.Now(),
			})
			if err != nil {
				websocket.JSON.Send(ws, gin.H{"error": err.Error()})
				return
			}
			server.broadcast(sessionId, msg)
		case "editObject":
			var payload dto.EditObjectPositionRequest
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				log.Println("Invalid payload for editObject:", err)
				continue
			}
			err := server.queries.UpsertObjectPosition(context.Background(), db.UpsertObjectPositionParams{
				SessionID:        payload.SessionID,
				ObjectID:         payload.ObjectID,
				XPos:             int32(payload.X),
				YPos:             int32(payload.Y),
				IsVisible:        *payload.IsVisible,
				ModificationDate: time.Now(),
			})
			if err != nil {
				websocket.JSON.Send(ws, gin.H{"error": err.Error()})
				return
			}
			server.broadcast(sessionId, msg)
		default:
			log.Println("Unknown message type:", msg.Type)
		}
	}
}

func (server *Server) broadcast(sessionId int64, msg dto.WebsocketMessage) {
	for ws := range server.sessionConns[sessionId] {
		go func(ws *websocket.Conn) {
			if err := websocket.JSON.Send(ws, msg); err != nil {
				log.Println("Send error:", err)
			}
		}(ws)
	}
}
