package main

import (
	"log"

	gomail "gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "JANGAN DIBALAS <noreply@suratupt.id>"
const CONFIG_AUTH_EMAIL = "emailsaya@gmail.com"
const CONFIG_AUTH_PASSWORD = "passwordemailservicesaya"

func main() {
	mailer := gomail.NewMessage()

	// Set email header
	mailer.SetHeader("From", CONFIG_SENDER_NAME)

	// Set email receiver
	mailer.SetHeader("To", "emailtarget@protonmail.com")

	// Set email subject
	mailer.SetHeader("Subject", "Giveaway")

	// Set email body
	mailer.SetBody("text/html", "<h2>Silakan KLIK tombol giveaway berikut ini <span>untuk menerima HADIAH!</span></h2> <br/> <a href=\"https://www.google.com\">HADIAH</a>")

	// Setting SMTP server
	dialer := gomail.NewDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatalf("[ERROR]: %v", err.Error())
	}

	log.Println("Mail sent!")
}
