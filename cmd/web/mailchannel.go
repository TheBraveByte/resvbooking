package main

import (
	"fmt"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// Using goroutine and channel to send mail to both the customer and the owner
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

	// to check for the email template
	if ml.MailTemplate == "" {
		email.SetBody(mail.TextHTML, ml.MailContent)
	} else {
		//this read the content in the email template html file
		emailData, err := ioutil.ReadFile(fmt.Sprintf("./email/%v", ml.MailTemplate))
		if err != nil {
			app.ErrorLog.Printf("Error reading mail template content ::", err)
		}
		//converted the byte data to a readable string
		emailText := string(emailData)
		emailText = strings.Replace(emailText, "[%body%]", ml.MailContent, 1)
		email.SetBody(mail.TextHTML, emailText)
	}
	err = email.Send(client)
	if err != nil {
		log.Println("error sending email")
	} else {
		log.Println("Email sent successfully")
	}

}
