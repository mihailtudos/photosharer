package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mihailtudos/photosharer/models"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "3ee29496783da1"
	password = "34731cfee5d6d5"
)

func main() {
	email := models.Email{
		From:      "test@renect.co.uk",
		To:        "mihairmcr7@gmail.com",
		Subject:   "Testing emails in Go",
		Plaintext: "This is the body of the email",
		HTML: `<h1>Hi there!</h1>
				<p>This is a test email so if you receive it it's great</p>
				</br></br> 
				<p>Regards</p>
				<p>Mihail</p>`,
	}

	es := models.NewEmailService(models.SMTPConfig{Host: host, Port: port, Username: username, Password: password})
	err := es.Send(email)

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
