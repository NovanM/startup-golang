package campaign

import "bwastartup/user"

type CampaignDetail struct {
	ID int `uri:"id" binding:"required"`
}
type CreateCampaign struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	User             user.User
}
type CampaignImageUpload struct {
	CampaignId int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       user.User
}
