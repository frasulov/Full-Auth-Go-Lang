package repo

import (
	"auth/models"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	connection *gorm.DB
}

func GetNewUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{
		connection: conn,
	}
}

func (userRepository *UserRepository) GetConnection() *gorm.DB {
	return userRepository.connection
}

func (userRepository *UserRepository) Init() {
	userRepository.connection.AutoMigrate(&models.User{})
}

func (userRepository *UserRepository) del(id uint) {
	userRepository.connection.Delete(&models.User{}, id)
}

func (userRepository *UserRepository) GetUserById(id string) (models.User, error) {
	var user models.User
	result := userRepository.connection.Model(&models.User{}).Where("id = ?", id).First(&user)
	return user, result.Error
}

func (userRepository *UserRepository) Save(user *models.User) error {
	return userRepository.connection.Create(&user).Error
}

func (userRepository *UserRepository) GetUserByMailOrUserName(username, email string) error {
	var result bool
	dbResult := userRepository.connection.Raw(`
			select
				case
					when count(*)=0 then false
					else true
				end
			from users where username = ? or email = ?;	
		`, username, email).Scan(&result)

	if dbResult.Error != nil || result {
		return fmt.Errorf("The user with the email/username is already exist!")
	}
	return nil
}

func (userRepository *UserRepository) GetUserByMail(email string) (models.User, error) {
	var user models.User
	if err := userRepository.connection.Where(models.User{Email: email}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (userRepository *UserRepository) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := userRepository.connection.Where(models.User{Email: email}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (userRepository *UserRepository) UpdateActive(id string) error {
	return userRepository.connection.Model(&models.User{}).Where("id = ?", id).Update("is_verified", true).Error
}

func (userRepository *UserRepository) SetPassword(id string, password string) error {
	return userRepository.connection.Model(&models.User{}).Where("id = ?", id).Update("password", password).Error
}
