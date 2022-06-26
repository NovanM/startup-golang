package transactions

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserId     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Campaign campaign.Campaign
}
