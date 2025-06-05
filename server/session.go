package server

import (
	"database/sql"
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/EliriaT/dnd-user-service/server/dto"
	"github.com/gin-gonic/gin"
	"net/http"
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
		IsVisible:        true,
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
