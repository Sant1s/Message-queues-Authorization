package main

import (
	"fmt"
	"github.com/Sant1s/MessageQueueTask/pkg/email"
	"github.com/Sant1s/MessageQueueTask/pkg/qeue"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("html", "index.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("html", "forgot_password.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func execSendPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("loginInput")
	mail := r.FormValue("emailInput")

	err := email.SendMessageLoginUser(login, mail)
	if err != nil {
		log.Println(err.Error())
	}

	path := filepath.Join("html", "index.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func execForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("loginInput")
	mail := r.FormValue("emailInput")

	err := email.SendMessageResetPassoword(login, mail)
	if err != nil {
		log.Println(err.Error())
	}

	path := filepath.Join("html", "index.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	params := u.Query()

	login := params.Get("login")
	mail := params.Get("email")
	token := params.Get("token")

	same, err := email.CheckTokenIsSame(login, mail, token, email.REGISTRATION_ACTON)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/invalid_token", http.StatusFound)
		return
	}

	if !same {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("html", "registration.html"))
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, struct {
		Login string
		Email string
		Token string
	}{
		Login: login,
		Email: mail,
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func execRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	params := u.Query()

	login := params.Get("login")
	mail := params.Get("email")

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	password := r.FormValue("passwordInput")
	// Вот тут надо дописать просто добавление пароля
	// и редирект на стартовую страницу
	// Там еще нужно будет добавить в хранении данных
}

func main() {
	go qeue.RunConsumer()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reset_password", resetPasswordHandler)
	http.HandleFunc("/exec/send_password_token", execSendPasswordTokenHandler)
	http.HandleFunc("/exec/forgot_password", execForgotPasswordHandler)

	http.HandleFunc("/registration", registrationHandler)
	http.HandleFunc("/exec/registration", execRegistrationHandler)

	http.HandleFunc("/invalid_token", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("html", "invalid_token.html")

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	})

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}

	select {}
}
