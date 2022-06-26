package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (ch *campaignHandler) GetCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := ch.service.GetCampaigns(userId)
	if err != nil {
		response := helper.APIResponse("Error get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formater := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)
}

func (ch *campaignHandler) GetCampaignDetail(c *gin.Context) {
	var input campaign.CampaignDetail

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get a campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := ch.service.GetCampaign(input.ID)
	if err != nil {
		response := helper.APIResponse("Failed to get a campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formater := campaign.FormatDetailCampaign(campaignDetail)
	response := helper.APIResponse("Detail of campaign", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)

}

func (ch *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaign
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create campaign failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	newCampaign, err := ch.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formater := campaign.FormatCampaign(newCampaign)
	response := helper.APIResponse("Successfuly created campaign", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)
	return
}

func (ch *campaignHandler) UpdatedCampaign(c *gin.Context) {
	var inputID campaign.CampaignDetail
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to updated a campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData campaign.CreateCampaign

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update campaign failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := ch.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formater := campaign.FormatCampaign(updatedCampaign)
	response := helper.APIResponse("Successfuly updated campaign", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)
	return
}

func (ch *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CampaignImageUpload
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Upload image campaign failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}	
	file, err := c.FormFile("file")
	if err != nil {

		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	input.User = currentUser
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {

		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = ch.service.CreateImageCampaign(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaing image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
