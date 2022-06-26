package transactions

import (
	"bwastartup/campaign"
	"bwastartup/payment"
	"errors"
	"strconv"
)

type Service interface {
	GetCampaignTransactionsByID(input TransactionCampaignInput) ([]Transaction, error)
	GetTransactionUserByID(userId int) ([]Transaction, error)
	CreateUserTransaction(input TransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNoficationInput) error
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepo campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepo, paymentService}
}

func (s *service) GetCampaignTransactionsByID(input TransactionCampaignInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not owned campaign")
	}

	transactions, err := s.repository.GetTransactionCampaignByID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *service) GetTransactionUserByID(userId int) ([]Transaction, error) {
	transaction, err := s.repository.FindTransactionByUserID(userId)
	if err != nil {
		return []Transaction{}, err
	}
	return transaction, nil
}

func (s *service) CreateUserTransaction(input TransactionInput) (Transaction, error) {
	var transaction Transaction = Transaction{
		CampaignID: input.CampaignID,
		Amount:     input.Amount,
		Status:     "pending",
	}
	transaction.UserId = input.User.ID
	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}
	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNoficationInput) error {
	trasaction_id, _ := strconv.Atoi(input.OrderID)
	trasaction, err := s.repository.FindByID(trasaction_id)
	if err != nil {
		return err
	}
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		trasaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		trasaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		trasaction.Status = "cancelled"
	}
	updatedTrasaction, err := s.repository.Update(trasaction)
	if err != nil {
		return err
	}
	campaign, err := s.campaignRepository.FindByID(updatedTrasaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTrasaction.Status == "paid" {
		campaign.BackerCount++
		campaign.CurrentAmount += updatedTrasaction.Amount
		_, err = s.campaignRepository.Updated(campaign)
		if err != nil {
			return err
		}

	}

	return nil

}
