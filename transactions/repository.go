package transactions

import "gorm.io/gorm"

type Repository interface {
	GetTransactionCampaignByID(campID int) ([]Transaction, error)
	FindTransactionByUserID(userId int) ([]Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
	FindByID(id int) (Transaction, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionCampaignByID(campID int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", campID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
func (r *repository) FindTransactionByUserID(userID int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transaction).Error

	if err != nil {
		return transaction, nil
	}
	return transaction, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
func (r *repository) FindByID(id int) (Transaction, error) {
	var transaction Transaction
	err := r.db.Where("id = ?", id).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
