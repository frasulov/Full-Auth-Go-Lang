package repo

import (
	"auth/models"
	"gorm.io/gorm"
)

type PasswordForgotRequestRepository struct {
	connection		*gorm.DB
}

func GetNewPasswordForgotRequestRepository(conn *gorm.DB) *PasswordForgotRequestRepository{
	return &PasswordForgotRequestRepository{
		connection: conn,
	}
}

func (passwordForgotRequestRepository * PasswordForgotRequestRepository) Init() {
	passwordForgotRequestRepository.connection.AutoMigrate(&models.PasswordForgotRequest{})
}

func (passwordForgotRequestRepository * PasswordForgotRequestRepository) FindById(id string) (models.PasswordForgotRequest, error) {
		var pfr models.PasswordForgotRequest
		err := passwordForgotRequestRepository.connection.Where("id = ? and is_active = ?", id, true).First(&pfr).Error
		return pfr, err
}

func (passwordForgotRequestRepository * PasswordForgotRequestRepository) UpdateActive(id string) error {
	return passwordForgotRequestRepository.connection.Model(&models.PasswordForgotRequest{}).Where("id = ?", id).Update("is_active", false).Error
}

func (passwordForgotRequestRepository * PasswordForgotRequestRepository) Save(request * models.PasswordForgotRequest) (*models.PasswordForgotRequest, error) {
	if err := passwordForgotRequestRepository.connection.Create(request).Error; err != nil{
		return request, err
	}
	return request, nil
}

