package tickets

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	TICKET_EMAIL_TEMPLATE    = "d-c706878f05bd4d22be6a5eb22b5a4dc8"
	ATTENDING_EMAIL_TEMPLATE = "d-e07fc7f75f2c4f7e9ffd97f48ae894d4"
)

func sendTicketEmail(apikey, email, name, ticketPath string) {
	from := mail.NewEmail("Socious", "info@socious.io")
	to := mail.NewEmail(name, email)

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(TICKET_EMAIL_TEMPLATE)

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetDynamicTemplateData("name", name)

	message.AddPersonalizations(personalization)

	// Attach a file
	fileBytes, err := os.ReadFile(ticketPath)
	if err != nil {
		log.Printf("Failed to read file: %v \n", err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	attachment := mail.NewAttachment()
	attachment.Content = encoded
	attachment.Type = "application/pdf"
	attachment.Filename = "ticket.pdf"
	attachment.Disposition = "attachment"

	message.AddAttachment(attachment)

	client := sendgrid.NewSendClient(apikey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Send error: %v \n", err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}

func sendAttendingEmail(apikey, email, name string) {
	from := mail.NewEmail("Socious", "info@socious.io")
	to := mail.NewEmail(name, email)

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(ATTENDING_EMAIL_TEMPLATE)

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.SetDynamicTemplateData("name", name)

	message.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(apikey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Send error: %v \n", err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}
