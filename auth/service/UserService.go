package service

import (
	"auth/models"
	"auth/repo"
	"fmt"
)

type UserService struct {
	repository repo.UserRepository
}

func GetNewService(userRepository repo.UserRepository) *UserService {
	return &UserService{
		repository: userRepository,
	}
}

func (userService * UserService) FindUserByEmailOrUsername(username string) (models.User,error){
	user, err := userService.repository.FindUserByEmailOrUsername(username)
	if err != nil{
		return user, fmt.Errorf("No such user with the email/username")
	}
	return user, nil
}

func (userService * UserService) GetUser(id uint) {
	user := userService.repository.GetUserById(id)
	fmt.Println(user)
}

func (userService *UserService) RegisterUser(user *models.User) error {
	err := userService.repository.GetUserByMailOrUserName(user.Username, user.Email)
	if err != nil {
		return err
	}
	err = userService.repository.Save(user)
	if err != nil{
		return fmt.Errorf("Error while creating new user: %v", err.Error())
	}
	return nil
}