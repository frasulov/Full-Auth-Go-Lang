package service

import (
	"auth/config"
	"auth/mail"
	"auth/models"
	"auth/repo"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"time"
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

func (userService * UserService) ResetPassword(password, confirmPassword, uuid string) ([]byte, error) {
	passwordForgotRequestRepository := repo.GetNewPasswordForgotRequestRepository(userService.repository.GetConnection())
	passwordForgotRequestRepository.Init()
	if len(password) < 8 {
		return nil, fmt.Errorf("Password length should be more than or equal to 8!")
	}
	if password != confirmPassword {
		return nil, fmt.Errorf("Passwords are not same!")
	}
	pfr, err := passwordForgotRequestRepository.FindById(uuid)
	if err != nil {
		return nil, fmt.Errorf("Token is not valid!")
	}
	if pfr.CreatedAt.Add(time.Hour*24).Unix() > time.Now().Unix() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil{
			return nil, fmt.Errorf("Unable to hash password")
		}
		err = userService.repository.SetPassword(pfr.UserId, string(hashedPassword))
		if err != nil {
			return nil, err
		}
		passwordForgotRequestRepository.UpdateActive(uuid)
		return []byte("Password has changes succesfully!"), nil
	}else {
		passwordForgotRequestRepository.UpdateActive(uuid)
		return nil, fmt.Errorf("Token has expired")
	}
}

func (userService * UserService) SendForgotPasswordMail(email string) ([]byte, error) {
	passwordForgotRequestRepository := repo.GetNewPasswordForgotRequestRepository(userService.repository.GetConnection())
	passwordForgotRequestRepository.Init()
	user, err := userService.repository.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("No such user with the email")
	}

	pfr := models.PasswordForgotRequest{
		UserId: user.ID,
	}
	passwordForgotRequestRepository.Save(&pfr)
	link := fmt.Sprintf("%s:%v/reset-password/%s", config.Configuration.Server.Host, config.Configuration.Server.Port, pfr.ID)
	auth := smtp.PlainAuth("", config.Configuration.Mail.Username, config.Configuration.Mail.Password, config.Configuration.Mail.Host)
	templateData := mail.NewTemplateData(user.FirstName + " " + user.LastName, link)
	r := mail.NewMailRequest([]string{user.Email}, "Forgot Password", "Reset Your Password Please")
	if err := r.ParseTemplate("./templates/forgot_password.html", templateData); err == nil {
		go func() {
			ok, err := r.SendEmail(auth)
			fmt.Println(ok)
			fmt.Println(err)
		}()
	}else{
		return nil, err
	}
	return []byte("Mail has been sent succesfully"), nil
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
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil{
		return fmt.Errorf("Unable to hash password")
	}
	user.Password = string(password)
	err = userService.repository.Save(user)
	if err != nil{
		return fmt.Errorf("Error while creating new user: %v", err.Error())
	}
	return nil
}

func (userService * UserService) ActivateUser(id uint) error{
	return userService.repository.UpdateActive(id)
}