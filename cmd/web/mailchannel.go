package main

import (
	"github.com/dev-ayaa/resvbooking/pkg/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

func mailRoutes() {
	go func() {
		for {
			mailMsg := <-app.MailChannel
			ListenToMailChannel(mailMsg)
		}
	}()
}

func ListenToMailChannel(ml models.MailData) {
	mailServer := mail.NewSMTPClient()
	mailServer.Host = "localhost"
	mailServer.Port = 1025
	mailServer.ConnectTimeout = 5 * time.Second
	mailServer.SendTimeout = 5 * time.Second
	mailServer.KeepAlive = false

	client, err := mailServer.Connect()
	if err != nil {
		log.Println("Error connecting to the mail server")
	}

	email := mail.NewMSG()
	email.SetFrom(ml.Sender).AddTo(ml.Receiver).SetSubject(ml.MailSubject)
	email.SetBody(mail.TextHTML, ml.MailContent)
	err = email.Send(client)
	if err != nil {
		log.Println("error sending email")
	} else {
		log.Println("Email sent successfully")
	}

}
