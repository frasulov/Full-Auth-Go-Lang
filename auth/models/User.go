package models

import (
	"auth/config"
	"auth/mail"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"net/smtp"
)

type User struct {
	Base
	IsVerified  bool           `json:"is_verified"`
	BirthDate   Date           `json:"birth_date"`
	FirstName   string         `json:"first_name"`
	Gender      sql.NullString `json:"gender"`
	LastName    string         `json:"last_name"`
	Username    sql.NullString `json:"username"`
	Email       string         `json:"email"`
	Country     sql.NullString `json:"country"`
	PhoneNumber sql.NullString `json:"phone_number"`
	Password    string         `json:"password"`
	PublicKey   sql.NullString `json:"public_key"`
	PrivateKey  sql.NullString `json:"private_key"`
	Role        sql.NullString `json:"role"`
	IsCompleted bool           `json:"is_completed"gorm:"default:false"`
}

func (user *User) AfterCreate(conn *gorm.DB) (err error) {
	userVerificationTokenRepository := GetNewUserVerificationTokenRepository(conn)
	fmt.Println(userVerificationTokenRepository)
	userVerificationToken := UserVerificationToken{
		UserId: user.ID,
	}
	userVerificationTokenRepository.Save(&userVerificationToken)
	link := fmt.Sprintf("%s:%v/confirm-email/%s", config.Configuration.Server.Host, config.Configuration.Server.Port, userVerificationToken.ID)
	auth := smtp.PlainAuth("", config.Configuration.Mail.Username, config.Configuration.Mail.Password, config.Configuration.Mail.Host)
	templateData := mail.NewTemplateData(user.FirstName+" "+user.LastName, link)
	r := mail.NewMailRequest([]string{user.Email}, "Mail Confirmation", "Confirm Your mail Please")
	if err := r.ParseTemplate("./templates/confirm_email.html", templateData); err == nil {
		go func() {
			ok, err := r.SendEmail(auth)
			fmt.Println(ok)
			fmt.Println(err)
		}()
	} else {
		return err
	}
	return nil
}
