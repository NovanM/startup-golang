package campaign

import (
	"strings"
)

type CampaignFormater struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `jsaon:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormater {
	formater := CampaignFormater{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
	}
	formater.ImageUrl = ""
	if len(campaign.CampaignImages) != 0 {
		formater.ImageUrl = campaign.CampaignImages[0].FileName
	}

	return formater
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormater {
	campaignsFormater := []CampaignFormater{}

	for _, v := range campaigns {
		campaignFormater := FormatCampaign(v)
		campaignsFormater = append(campaignsFormater, campaignFormater)
	}
	return campaignsFormater
}

type CampaignDetailFormater struct {
	ID               int                      `json:"id"`
	UserID           int                      `json:"user_id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	ImageUrl         string                   `json:"image_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	Slug             string                   `jsaon:"slug"`
	Description      string                   `json:"description"`
	BeckerCount      int                      `json:"becker_count"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormater     `json:"user"`
	Images           []CampaignImagesFormater `json:"images"`
}

type CampaignUserFormater struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CampaignImagesFormater struct {
	ImageUrl  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatDetailCampaign(campaign Campaign) CampaignDetailFormater {
	formater := CampaignDetailFormater{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Description:      campaign.Description,
		Slug:             campaign.Slug,
		BeckerCount:      campaign.BackerCount,
	}
	formater.ImageUrl = ""
	if len(campaign.CampaignImages) != 0 {
		formater.ImageUrl = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}
	formater.Perks = perks

	user := campaign.User
	formater.User.Name = user.Name
	formater.User.ImageUrl = user.AvatarFileName

	campaignImages := []CampaignImagesFormater{}
	for _, image := range campaign.CampaignImages {
		campaignImage := CampaignImagesFormater{
			ImageUrl: image.FileName,
		}
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImage.IsPrimary = isPrimary
		campaignImages = append(campaignImages, campaignImage)
	}
	formater.Images = campaignImages

	return formater

}
