package tickets

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"socious/src/apps/models"
	"socious/src/config"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/socious-io/goaccount"
	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"
	"github.com/stripe/stripe-go/v81"
)

var (
	configPath          = flag.String("c", "config.yml", "Path to the configuration file")
	mode                = flag.String("m", "", "Operation mode: producer, customer-consumer, email-consumer, pdf-consumer, publish-customer")
	ticketPath          = flag.String("t", "", "Path to ticket template")
	ticketsGeneratedDir = flag.String("o", "", "Directory of tickets")
	sendgridApiKey      = flag.String("ak", "", "Sendgrid api key")

	// this use on publish-customer
	email      = flag.String("email", "", "Email to send")
	name       = flag.String("name", "", "Name to send")
	company    = flag.String("company", "", "compay name to send")
	ticketType = flag.String("type", "", "ticket type to send")

	event *models.Event
	nc    *nats.Conn
)

const (
	PUBLISH  = "publish-customer"
	CUSTOMER = "customer-consumer"
	EMAIL    = "email-consumer"
	PDF      = "pdf-consumer"

	profileAddress = "https://app.socious.io/profile/users/%s/view"
)

type Customer struct {
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Company    string    `json:"company"`
	TicketType string    `json:"ticket_type"`
}

func Run() {
	flag.Parse()
	config.Init(*configPath)
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "",
	})
	gomq.Connect()

	if c, err := nats.Connect(config.Config.Nats.Url, nats.Token(config.Config.Nats.Token)); err == nil {
		nc = c
	} else {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	defer nc.Close()

	//Configure Socious ID SDK
	goaccount.Setup(config.Config.GoAccounts)
	// Set your secret key
	stripe.Key = config.Config.Payment.Fiats[0].ApiKey

	e, err := models.GetActiveEvent()
	if err != nil {
		log.Fatalf("there is no active event: %v", err)
		return
	} else {
		event = e
	}

	switch *mode {
	case PUBLISH:
		publishCustomer()
	case CUSTOMER:
		customerConsumer()
	case PDF:
		pdfConsumer()
	case EMAIL:
		emailConsumer()
	default:
		fetchPaymentLinks()
	}

}

func publishCustomer() {
	if name == nil || email == nil {
		log.Fatal("email and name are required")
		return
	}

	publish(consumerTitle(CUSTOMER), Customer{
		Name:       *name,
		Email:      *email,
		Company:    *company,
		TicketType: *ticketType,
	})
}

func customerConsumer() {
	_, err := nc.Subscribe(consumerTitle(CUSTOMER), func(msg *nats.Msg) {
		// Parse the message (format: "type|content")
		log.Printf("%s got %s \n", consumerTitle(CUSTOMER), string(msg.Data))
		cus := new(Customer)
		if err := json.Unmarshal(msg.Data, &cus); err != nil {
			log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
		}

		user, err := models.GetUserByEmail(cus.Email)
		if err != nil {
			names := strings.Split(cus.Name, " ")
			var (
				fName = names[0]
				lName = ""
			)
			if len(names) > 1 {
				lName = names[1]
			}

			accountUser := &goaccount.User{
				Email:     cus.Email,
				FirstName: &fName,
				LastName:  &lName,
			}
			if err := accountUser.Create(); err != nil {
				log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
				return
			}

			user = &models.User{
				ID:        accountUser.ID,
				FirstName: accountUser.FirstName,
				LastName:  accountUser.LastName,
				Email:     accountUser.Email,
				Username:  accountUser.Username,
				Events:    pq.StringArray{event.ID.String()},
			}
			if err := user.Upsert(context.Background()); err != nil {
				log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
				return
			}

		} else {
			if existsOnEvent(user.Events) {
				return
			}

			user.Events = append(user.Events, event.ID.String())
			if err := user.Upsert(context.Background()); err != nil {
				log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
				return
			}
		}

		cus.UserID = user.ID
		cus.Username = user.Username

		publish(consumerTitle(PDF), cus)

	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever
}

func pdfConsumer() {
	_, err := nc.Subscribe(consumerTitle(PDF), func(msg *nats.Msg) {
		// Parse the message (format: "type|content")
		log.Printf("%s got %s \n", consumerTitle(PDF), string(msg.Data))
		cus := new(Customer)
		if err := json.Unmarshal(msg.Data, cus); err != nil {
			log.Printf("Error on consumer pdf: %v | data: %s ", err, string(msg.Data))
		}
		if PdfGenerator(
			*ticketPath,
			fmt.Sprintf("%s/%s.pdf", *ticketsGeneratedDir, cus.Username),
			cus.Name,
			fmt.Sprintf(profileAddress, cus.Username),
			cus.Company,
			cus.TicketType,
		) {
			publish(consumerTitle(EMAIL), cus)
		}

	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever

}

func emailConsumer() {
	_, err := nc.Subscribe(consumerTitle(EMAIL), func(msg *nats.Msg) {
		log.Printf("%s got %s \n", consumerTitle(PDF), string(msg.Data))
		cus := new(Customer)
		if err := json.Unmarshal(msg.Data, cus); err != nil {
			log.Printf("Error on consumer pdf: %v | data: %s ", err, string(msg.Data))
		}
		ticket := fmt.Sprintf("%s/%s.pdf", *ticketsGeneratedDir, cus.Username)

		apiKey := config.Config.SendgridApiKey
		if *sendgridApiKey != "" {
			apiKey = *sendgridApiKey
		}

		sendEmail(apiKey, cus.Email, cus.Name, ticket)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever
}

func existsOnEvent(events pq.StringArray) bool {
	for _, e := range events {
		if e == event.ID.String() {
			return true
		}
	}
	return false
}

func publish(channel string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return nc.Publish(channel, body)
}

func consumerTitle(channel string) string {
	return fmt.Sprintf("ticketing-%s", channel)
}
