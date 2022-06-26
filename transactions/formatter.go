package transactions

import (
	"time"
)

type TransactionFormater struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionPaymentFormatter struct {
	ID         int    `json:"id"`
	CampaignId int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

type UserTransactionFormater struct {
	ID        int              `json:"id"`
	Amount    int              `json:"amount"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	Campaign  CampaignFormater `json:"campaign"`
}

type CampaignFormater struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func FormaterTransaction(transaction Transaction) TransactionFormater {
	return TransactionFormater{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

}

func FormaterTransactions(Transaction []Transaction) []TransactionFormater {
	if len(Transaction) == 0 {
		return []TransactionFormater{}
	}

	var transactions []TransactionFormater
	for _, v := range Transaction {
		formatter := FormaterTransaction(v)
		transactions = append(transactions, formatter)
	}

	return transactions
}

func FormatterUserTransaction(transaction Transaction) UserTransactionFormater {
	var campaignTransaction CampaignFormater
	campaignTransaction.Name = transaction.Campaign.Name
	campaignTransaction.ImageUrl = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignTransaction.ImageUrl = transaction.Campaign.CampaignImages[0].FileName
	}

	return UserTransactionFormater{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaignTransaction,
	}
}

func FormatterUserTransactions(transaction []Transaction) []UserTransactionFormater {
	if len(transaction) == 0 {
		return []UserTransactionFormater{}
	}
	var userTransactions []UserTransactionFormater
	for _, v := range transaction {
		userTransaction := FormatterUserTransaction(v)
		userTransactions = append(userTransactions, userTransaction)
	}
	return userTransactions
}

func FormaterPaymentTransaction(transaction Transaction) TransactionPaymentFormatter {
	return TransactionPaymentFormatter{
		ID:         transaction.ID,
		CampaignId: transaction.CampaignID,
		Code:       transaction.Code,
		Status:     transaction.Status,
		UserID:     transaction.UserId,
		Amount:     transaction.Amount,
		PaymentURL: transaction.PaymentURL,
	}
}
