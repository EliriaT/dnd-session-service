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
