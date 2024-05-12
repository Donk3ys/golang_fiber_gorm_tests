package notification

import (
	"api/src/constants"
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/gofiber/fiber/v2/log"
)

type Email struct {
	Client IEmailClient
}

func (e *Email) SendUserVerificationCode(email string, code uint16) error {
	to := []string{email}
	href := fmt.Sprintf("%s/verify?email=%s&code=%d", constants.WEBSITE_URL, email, code)
	subject := "Sign Up Email Verification"
	templatePath := "templates/signup_code.html"
	data := struct {
		Href string
	}{
		Href: href,
	}

	return e.Client.SendMail("", to, subject, templatePath, data)
}

func (e *Email) SendPasswordResetCode(email string, code uint16) error {
	to := []string{email}
	subject := fmt.Sprintf("Password Reset Code")
	templatePath := "templates/password_reset_code.html"
	data := struct {
		Code uint16
	}{
		Code: code,
	}

	return e.Client.SendMail("", to, subject, templatePath, data)
}

func (e *Email) SendPasswordResetSuccess(email string) error {
	to := []string{email}
	subject := "Password Reset Success"
	templatePath := "templates/password_reset_success.html"

	return e.Client.SendMail("", to, subject, templatePath, nil)
}

// Interface for test mock
type IEmailClient interface {
	SendMail(from string, to []string, subject string, templatePath string, data interface{}) error
}
type SMTP struct {
	User     string
	Password string
	Host     string
	Port     string
}

func (s SMTP) SendMail(from string, to []string, subject string, templatePath string, data interface{}) error {
	if from == "" {
		from = "no-reply@support.com"
	}

	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("From: %s\nSubject: %s\n%s\n\n", from, subject, headers)))

	err = template.Execute(&body, data)
	if err != nil {
		return err
	}

	// Sending email.
	err = smtp.SendMail(s.Host+":"+s.Port, auth, from, to, body.Bytes())
	if err != nil {
		log.Errorf("[Email Error]: %s", err)
		return err
	}
	return nil
}
