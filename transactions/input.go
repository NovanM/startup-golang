package transactions

import "bwastartup/user"

type TransactionCampaignInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

type TransactionInput struct {
	CampaignID int `json:"campaign_id"`
	Amount     int `json:"amount"`

	User user.User
}

type TransactionNoficationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
