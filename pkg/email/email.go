package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	SMTP_SERVER_NAME      = "smtp.yandex.ru"
	SMTP_SERVER_PORT      = 465
	SENDER_EMAIL          = "nam.i@phystech.edu"
	SENDER_EMAIL_PASSWORD = "06051999qew"

	REGISTRATION_ACTON    = "registration"
	RESET_PASSWORD_ACTION = "reset_password"

	SERVER_ADDRES = "localhost:8080"
)

func generateLink(login, email, action string) (string, error) {
	userData := clients[Clients{login, email}]
	var token string

	if action == REGISTRATION_ACTON {
		token = userData.userRegistrationToken.token
	} else {
		token = userData.userResetPasswordToken.token
	}

	return fmt.Sprintf("http://%s/%s?login=%s&email=%s&token=%s", SERVER_ADDRES, action, login, email, token), nil
}

func generateMailText(login, email, action string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	tmpl, err := template.ParseFiles("html/registration_tmpl.html")
	if err != nil {
		return nil, err
	}

	link, err := generateLink(login, email, action)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&buf, struct {
		Username   string
		Link       string
		ExpiryDate string
	}{
		Username:   login,
		Link:       link,
		ExpiryDate: time.Now().Add(48 * time.Hour).Format("02.01.2006 15:04:05"), // ссылка истекает через 48 часов
	})

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func SendMessageLoginUser(login, email string) error {
	err := updateRegistrationToken(login, email)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(SMTP_SERVER_NAME, SMTP_SERVER_PORT, SENDER_EMAIL, SENDER_EMAIL_PASSWORD)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mail, err := generateMailText(login, email, REGISTRATION_ACTON)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", SENDER_EMAIL)
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("Регистрация нового пользователя %s!", login))
	m.SetBody("text/html", mail.String())

	if err = dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendMessageResetPassoword(login, email string) error {
	err := updateResetPasswordToken(login, email)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(SMTP_SERVER_NAME, SMTP_SERVER_PORT, SENDER_EMAIL, SENDER_EMAIL_PASSWORD)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mail, err := generateMailText(login, email, RESET_PASSWORD_ACTION)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", SENDER_EMAIL)
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("Изменение пароля для %s", login))
	m.SetBody("text/html", mail.String())

	if err = dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
