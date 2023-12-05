package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-mail/mail/v2"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "3ee29496783da1"
	password = "34731cfee5d6d5"
)

func main() {
	//email := models.Email{}
	from := "test@renect.co.uk"
	to := "mihairmcr7@gmail.com"
	subject := "Testing emails in Go"
	plainText := "This is the body of the email"
	html := `<h1>Hi there!</h1>
				<p>This is a test email so if you receive it it's great</p>
				</br></br> 
				<p>Regards</p>
				<p>Mihail</p>
	`
	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainText)
	msg.AddAlternative("text/html", html)

	dialer := mail.NewDialer(host, port, username, password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent")
}

func passHasingExample() {
	secret := "MySuperSecretPhrase"
	password := "ThisIsMyPassword"

	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(password))

	output := h.Sum(nil)

	fmt.Println(hex.EncodeToString(output))
	// => 827e2efb277ea22df6e9559ecd5dd5448b7da6f1ba3d63fde0f14b91714e9bb7
}
