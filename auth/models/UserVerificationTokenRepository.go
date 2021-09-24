package models

import (
	"gorm.io/gorm"
)

type UserVerificationTokenRepository struct {
	connection *gorm.DB
}

func GetNewUserVerificationTokenRepository(conn *gorm.DB) *UserVerificationTokenRepository {
	return &UserVerificationTokenRepository{
		connection: conn,
	}
}

func (userVerificationTokenRepository *UserVerificationTokenRepository) Init() {
	userVerificationTokenRepository.connection.AutoMigrate(&UserVerificationToken{})
}

func (userVerificationTokenRepository *UserVerificationTokenRepository) Save(userVerificationToken *UserVerificationToken) error {
	result := userVerificationTokenRepository.connection.Create(&userVerificationToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (userVerificationTokenRepository *UserVerificationTokenRepository) FindByToken(token string) (UserVerificationToken, error) {
	var userVerificationToken UserVerificationToken
	result := userVerificationTokenRepository.connection.Where("id = ?", token).First(&userVerificationToken)
	if result.Error != nil {
		return userVerificationToken, result.Error
	}
	return userVerificationToken, nil
}
