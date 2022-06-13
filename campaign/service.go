package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaign(ID int) (Campaign, error)
	CreateCampaign(campaign CreateCampaign) (Campaign, error)
	UpdateCampaign(inputID CampaignDetail, campaign CreateCampaign) (Campaign, error)
	CreateImageCampaign(campaignImage CampaignImageUpload, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userId int) ([]Campaign, error) {

	if userId != 0 {
		campaigns, err := s.repository.FindByUserId(userId)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetCampaign(ID int) (Campaign, error) {
	campaign, err := s.repository.FindByID(ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaign) (Campaign, error) {
	var campaign = Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Perks:            input.Perks,
		GoalAmount:       input.GoalAmount,
		UserID:           input.User.ID,
	}
	slugURL := fmt.Sprintf("%s %v", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugURL)
	newCamapign, err := s.repository.Save(campaign)
	if err != nil {
		return newCamapign, err
	}
	return newCamapign, nil

}

func (s *service) UpdateCampaign(inputID CampaignDetail, inputData CreateCampaign) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}
	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not owned a campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Updated(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

func (s *service) CreateImageCampaign(campaignImage CampaignImageUpload, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindByID(campaignImage.CampaignId)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != campaignImage.User.ID {
		return CampaignImage{}, errors.New("Not owned a campaign")
	}

	isPrimay := 0
	if campaignImage.IsPrimary {
		isPrimay = 1
		_, err := s.repository.MarkAllIsprimary(campaignImage.CampaignId)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	imageCampaign := CampaignImage{
		CampaignID: campaignImage.CampaignId,
		FileName:   fileLocation,
		IsPrimary:  isPrimay,
	}

	newCampaign, err := s.repository.CreateImage(imageCampaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil

}
