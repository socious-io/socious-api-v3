package main

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
	"github.com/nats-io/nats.go"
	"github.com/socious-io/goaccount"
	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"

	"github.com/stripe/stripe-go/v81"
)

var (
	configPath          = flag.String("c", "config.yml", "Path to the configuration file")
	mode                = flag.String("m", "", "Operation mode: producer, customer-consumer, email-consumer, pdf-consumer")
	ticketPath          = flag.String("t", "", "Path to ticket template")
	ticketsGeneratedDir = flag.String("o", "", "Directory of tickets")
	event               *models.Event
	nc                  *nats.Conn
)

const (
	CUSTOMER = "customer-consumer"
	EMAIL    = "email-consumer"
	PDF      = "pdf-consumer"

	profileAddress = "https://app.socious.io/profile/users/%s/view"
)

type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
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

	if e, err := models.GetActiveEvent(); err == nil {
		event = e
	} else {
		log.Fatal("there is no active event")
		return
	}

	switch *mode {
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

func customerConsumer() {
	_, err := nc.Subscribe(consumerTitle(CUSTOMER), func(msg *nats.Msg) {
		// Parse the message (format: "type|content")
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
				Events:    []uuid.UUID{event.ID},
			}
			if err := user.Upsert(context.Background()); err != nil {
				log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
				return
			}

		} else {
			if existsOnEvent(user.Events) {
				return
			}

			user.Events = append(user.Events, event.ID)
			if err := user.Upsert(context.Background()); err != nil {
				log.Printf("Error on consumer customer: %v | data: %s ", err, string(msg.Data))
				return
			}
		}

		publish(consumerTitle(PDF), user)

	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever
}

func pdfConsumer() {
	_, err := nc.Subscribe(consumerTitle(PDF), func(msg *nats.Msg) {
		// Parse the message (format: "type|content")
		user := new(models.User)
		if err := json.Unmarshal(msg.Data, user); err != nil {
			log.Printf("Error on consumer pdf: %v | data: %s ", err, string(msg.Data))
		}
		if pdfGenerator(
			*ticketPath,
			fmt.Sprintf("%s/%s.pdf", *ticketsGeneratedDir, user.Username),
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			fmt.Sprintf(profileAddress, user.Username),
		) {
			publish(consumerTitle(EMAIL), user)
		}

	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever

}

func emailConsumer() {
	_, err := nc.Subscribe(consumerTitle(EMAIL), func(msg *nats.Msg) {
		user := new(models.User)
		if err := json.Unmarshal(msg.Data, user); err != nil {
			log.Printf("Error on consumer pdf: %v | data: %s ", err, string(msg.Data))
		}
		ticket := fmt.Sprintf("%s/%s.pdf", *ticketsGeneratedDir, user.Username)
		sendEmail(user.Email, fmt.Sprintf("%s %s", user.FirstName, user.LastName), ticket)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {} // block forever
}

func existsOnEvent(events []uuid.UUID) bool {
	for _, e := range events {
		if e == event.ID {
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
