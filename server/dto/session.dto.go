package dto

type GetSessionsByCampaignRequest struct {
	UserId     int64 `json:"userId" binding:"required"`
	CampaignId int64 `json:"campaignId" binding:"required"`
}
