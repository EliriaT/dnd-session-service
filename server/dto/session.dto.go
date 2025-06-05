package dto

type GetSessionsByCampaignRequest struct {
	UserId     int64 `json:"userId" binding:"required"`
	CampaignId int64 `json:"campaignId" binding:"required"`
}

type CreateSessionRequest struct {
	Name         string  `json:"name" binding:"required"`
	CampaignID   int64   `json:"campaignId" binding:"required"`
	MapID        int64   `json:"mapId" binding:"required"`
	AllowedChars []int64 `json:"allowedCharacters" binding:"required,min=1,dive,required"`
}

type EditCharacterPositionRequest struct {
	SessionID   int64 `json:"sessionId" binding:"required"`
	CharacterID int64 `json:"characterId" binding:"required"`
	X           int   `json:"x" binding:"required"`
	Y           int   `json:"y" binding:"required"`
}

type EditObjectPositionRequest struct {
	SessionID int64 `json:"sessionId" binding:"required"`
	ObjectID  int64 `json:"objectId" binding:"required"`
	X         int   `json:"x" binding:"required"`
	Y         int   `json:"y" binding:"required"`
}

type GetIdRequest struct {
	SessionID int64 `uri:"sessionId" binding:"required,min=1"`
}
