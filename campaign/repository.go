package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserId(userId int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Updated(campaign Campaign) (Campaign, error)
	CreateImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllIsprimary(campaignID int) (bool, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.DB.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserId(userId int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.DB.Where("user_id = ?", userId).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.DB.Where("id = ?", ID).Preload("User").Preload("CampaignImages").Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}
func (r *repository) Save(campaign Campaign) (Campaign, error) {

	err := r.DB.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

func (r *repository) Updated(campaign Campaign) (Campaign, error) {
	err := r.DB.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := r.DB.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *repository) MarkAllIsprimary(campaignID int) (bool, error) {
	err := r.DB.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
