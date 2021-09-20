package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type MailRequest struct {
	from    string
	to      []string
	subject string
	body    string
}

type TemplateData struct {
	Name string
	Url  string
}

func NewTemplateData(name, url string) *TemplateData {
	return &TemplateData{
		Url: url,
		Name: name,
	}
}

func NewMailRequest(to []string, subject, body string) *MailRequest {
	return &MailRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *MailRequest) SendEmail(auth smtp.Auth) (bool, error) {
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + MIME + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "dhanush@geektrust.in", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *MailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
