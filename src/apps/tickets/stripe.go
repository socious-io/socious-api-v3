package tickets

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/paymentlink"
)

func fetchPaymentLinks() {
	params := &stripe.PaymentLinkListParams{
		Active: stripe.Bool(true),
	}

	// Optional: Set pagination parameters
	params.Limit = stripe.Int64(10)

	// List active payment links
	i := paymentlink.List(params)
	for i.Next() {
		link := i.PaymentLink()
		fmt.Printf("Payment Link ID: %s, URL: %s\n", link.ID, link.URL)
		fetchSuccessfulPaymentsForLink(link.ID)
	}

	if err := i.Err(); err != nil {
		log.Printf("Error listing payment links: %v\n", err)
	}
}

func fetchSuccessfulPaymentsForLink(paymentLinkID string) {

	// Create parameters for listing checkout sessions
	params := &stripe.CheckoutSessionListParams{}
	params.Filters.AddFilter("payment_link", "", paymentLinkID)
	params.Filters.AddFilter("status", "", "complete")

	// Optional: Add expansion to get more customer details
	params.AddExpand("data.customer")
	params.AddExpand("data.customer_details")

	// List checkout sessions for the payment link
	i := session.List(params)
	for i.Next() {
		s := i.CheckoutSession()
		fmt.Printf("Session ID: %s\n", s.ID)
		customer := new(Customer)
		// Get customer information
		if s.CustomerEmail != "" {
			fmt.Printf("Customer Email: %s\n", s.CustomerEmail)
		}

		if s.CustomerDetails != nil {
			customer.Name = s.CustomerDetails.Name
			customer.Email = s.CustomerDetails.Email
		}

		if s.Customer != nil {
			customer.Name = s.Customer.Name
			customer.Email = s.Customer.Email
		}

		for _, field := range s.CustomFields {
			if field.Key == "companyname" {
				customer.Company = field.Text.Value
			}
		}

		customer.TicketType = linkType(paymentLinkID)

		publish(consumerTitle(CUSTOMER), customer)

	}

	if err := i.Err(); err != nil {
		log.Printf("Error listing checkout sessions: %v\n", err)
	}
}

func linkType(linkID string) string {
	switch linkID {
	case "plink_1RsdPvFiHSKRe5D1sErI3vNO":
		return "Standard"
	case "plink_1RkiVZFiHSKRe5D1enUxMBhK":
		return "Corporate"
	case "plink_1RkglwFiHSKRe5D1tKLALbNS":
		return "Senior"
	case "plink_1RgF5QFiHSKRe5D1gIEPi4Xv":
		return "Investor"
	case "plink_1RgF4QFiHSKRe5D1jkzfs5Uc":
		return "Startup"
	case "plink_1RgEoQFiHSKRe5D1hj1Pcv7i":
		return "Student"
	case "plink_1QnC2HFiHSKRe5D1V2GmnEdd":
		return "VIP"
	default:
		return "Standard"
	}
}
