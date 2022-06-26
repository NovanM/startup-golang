package handler

import (
	"bwastartup/helper"
	"bwastartup/transactions"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transactions.Service
}

func NewTTransactionHandler(service transactions.Service) *transactionHandler {
	return &transactionHandler{service}
}
func (t *transactionHandler) GetTransactionCampaign(c *gin.Context) {
	var input transactions.TransactionCampaignInput

	err := c.ShouldBindUri(&input)

	if err != nil {

		response := helper.APIResponse("Failed to get campaign transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	transaction, err := t.service.GetCampaignTransactionsByID(input)

	if err != nil {

		response := helper.APIResponse("Failed to get campaign transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := transactions.FormaterTransactions(transaction)
	response := helper.APIResponse("Detail of transaction", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (t *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	userTransactions, err := t.service.GetTransactionUserByID(userId)

	if err != nil {

		response := helper.APIResponse("Failed to get transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := transactions.FormatterUserTransactions(userTransactions)
	response := helper.APIResponse("Detail of transaction", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (t *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transactions.TransactionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {

		response := helper.APIResponse("Failed to transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	newTransaction, err := t.service.CreateUserTransaction(input)
	if err != nil {

		response := helper.APIResponse("Failed to transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Detail of transaction", http.StatusOK, "success", transactions.FormaterPaymentTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (t *transactionHandler) GetNotification(c *gin.Context) {
	var input transactions.TransactionNoficationInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = t.service.ProcessPayment(input)
	if err != nil {
		response := helper.APIResponse("Failed process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
