package models

import (
	"auth/config"
	"auth/mail"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/smtp"
	"time"
)

type User struct {
	ID 				uint			`json:"id"`
	FirstName		string			`gorm:"size:32"json:"first_name"`
	LastName		string			`gorm:"size:32"json:"last_name"`
	Username 		string			`gorm:"size:64;not null;unique"json:"username"`
	Email			string			`gorm:"size:128;not null;unique"json:"email"`
	Password		string
	IsActive 		bool			`gorm:"default:false"json:"is_active"`
	CreatedAt 		time.Time		`gorm:"default:now()"json:"created_at"`
	UpdatedAt 		time.Time		`gorm:"default:now()"json:"updated_at"`
}

func (user * User) AfterCreate(conn *gorm.DB) (err error) {
	userVerificationTokenRepository := GetNewUserVerificationTokenRepository(conn)
	userVerificationTokenRepository.Init()
	token := uuid.New().String()
	userVerificationToken := UserVerificationToken{
		Token: token,
		UserId: user.ID,
		ExpiredAt: time.Now().Add(time.Minute*time.Duration(config.Configuration.Password.ActivateAccountTokenExpire)).Unix(),
	}
	userVerificationTokenRepository.Save(&userVerificationToken)
	link := fmt.Sprintf("%s:%v/confirm-email/%s", config.Configuration.Server.Host, config.Configuration.Server.Port, token)
	auth := smtp.PlainAuth("", config.Configuration.Mail.Username, config.Configuration.Mail.Password, config.Configuration.Mail.Host)
	templateData := mail.NewTemplateData(user.FirstName + " " + user.LastName, link)
	r := mail.NewMailRequest([]string{user.Email}, "Mail Confirmation", "Confirm Your mail Please")
	if err := r.ParseTemplate("./templates/confirm_email.html", templateData); err == nil {
		go func() {
			ok, err := r.SendEmail(auth)
			fmt.Println(ok)
			fmt.Println(err)
		}()
	}else{
		return err
	}
	return nil
}