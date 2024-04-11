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
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	params := u.Query()

	login := params.Get("login")
	mail := params.Get("email")
	token := params.Get("token")

	same, err := email.CheckTokenIsSame(login, mail, token, email.RESET_PASSWORD_ACTION)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/invalid_token", http.StatusNotFound)
		return
	}

	if !same {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	err = qeue.SendMessage(fmt.Sprintf("login=%s&email=%s&password=%s", login, mail, r.FormValue("passwordInput")), "password")
	if err != nil {
		log.Println(err.Error())
	}

	path := filepath.Join("html", "reset_password.html")

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
		http.Redirect(w, r, "/invalid_token", http.StatusNotFound)
		return
	}

	if !same {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
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

func sendRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("loginInput")
	mail := r.FormValue("emailInput")

	err := qeue.SendMessage(fmt.Sprintf("login=%s&email=%s", login, mail), email.REGISTRATION_ACTON)
	if err != nil {
		log.Println(err.Error())
	}

	path := filepath.Join("html", "send_register.html")

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

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
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

func sendForgotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("loginInput")
	mail := r.FormValue("emailInput")

	err := qeue.SendMessage(fmt.Sprintf("login=%s&email=%s", login, mail), email.RESET_PASSWORD_ACTION)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/invalid_token", http.StatusNotFound)
	}

	path := filepath.Join("html", "send_forgot.html")

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

func main() {
	go qeue.RunConsumer()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/send_register", sendRegisterHandler)
	http.HandleFunc("/forgot_password", forgotPasswordHandler)
	http.HandleFunc("/send_forgot", sendForgotHandler)
	http.HandleFunc("/registration", registrationHandler)
	http.HandleFunc("/reset_password", resetPasswordHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}

	select {}
}
