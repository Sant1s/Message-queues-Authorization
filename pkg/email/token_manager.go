package email

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

type UnknownUser struct {
	message string
}

func (u *UnknownUser) Error() string {
	return u.message
}

type ValidLink struct {
	token      string
	expireDate time.Time
}

type Links struct {
	userRegistrationToken  *ValidLink
	userResetPasswordToken *ValidLink
	userPassword           string
}

type Clients struct {
	userLogin string
	userEmail string
}

var clients map[Clients]Links

func generateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

func updateResetPasswordToken(login, email string) error {
	if val, ok := clients[Clients{login, email}]; ok {
		if time.Since(val.userResetPasswordToken.expireDate) < 48*time.Hour {
			token, err := generateToken()
			if err != nil {
				return err
			}
			val.userResetPasswordToken.token = token
			val.userRegistrationToken.expireDate = time.Now()
		}
	} else {
		return &UnknownUser{"Unknown user. Can not update token"}
	}
	return nil
}

func updateRegistrationToken(login, email string) error {
	if clients == nil {
		clients = make(map[Clients]Links)
	}

	if val, ok := clients[Clients{login, email}]; ok {
		if time.Since(val.userResetPasswordToken.expireDate) < 48*time.Hour {
			token, err := generateToken()
			if err != nil {
				return err
			}
			val.userResetPasswordToken.token = token
			val.userRegistrationToken.expireDate = time.Now()
		}
	} else {
		registerToken, err := generateToken()
		if err != nil {
			return nil
		}

		resetToken, err := generateToken()
		if err != nil {
			return err
		}

		clients[Clients{login, email}] = Links{
			userRegistrationToken: &ValidLink{
				token:      registerToken,
				expireDate: time.Now(),
			},
			userResetPasswordToken: &ValidLink{
				token:      resetToken,
				expireDate: time.Now(),
			},
		}
	}
	return nil
}

func CheckTokenIsSame(login, email, token, action string) (bool, error) {
	if client, ok := clients[Clients{login, email}]; ok {
		switch action {
		case REGISTRATION_ACTON:
			return client.userRegistrationToken.token == token, nil
		case RESET_PASSWORD_ACTION:
			return client.userResetPasswordToken.token == token, nil
		}
	} else {
		return false, &UnknownUser{"Unknown user. Can not check token"}
	}
	return true, nil
}

func GetToken(link string) string {
	parts := strings.Split(link, "token=")
	return parts[1]
}

func SetPassword(login, email, password string) error {
	if client, ok := clients[Clients{login, email}]; ok {
		client.userPassword = password
	} else {
		return &UnknownUser{"Unknown user. Can not set password"}
	}
	return nil
}
