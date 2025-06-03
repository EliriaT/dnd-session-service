package server

import (
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/EliriaT/dnd-user-service/server/dto"
	"github.com/gin-gonic/gin"
	"net/http"
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
