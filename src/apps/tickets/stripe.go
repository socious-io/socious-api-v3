package tickets

import (
	"fmt"
	"log"
	"strings"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/paymentlink"
)

func fetchPaymentLinks() {
	params := &stripe.PaymentLinkListParams{}
	params.Limit = stripe.Int64(500)
	params.AddExpand("data.after_completion")
	params.AddExpand("data.line_items")
	params.AddExpand("data.application")
	params.AddExpand("data.on_behalf_of")

	i := paymentlink.List(params)
	for i.Next() {
		link := i.PaymentLink()
		is2025 := false
		desc := link.LineItems.Data[0].Description

		if strings.Contains(desc, "2025") {
			is2025 = true
		}
		fmt.Printf("\n\n\n\n")
		fmt.Printf("\n\n\n\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("Payment Link ID: %s, URL: %s\n", link.ID, link.URL)
		fmt.Printf("Link Description: %s --  status: %t --  detected as 2025: %t\n", desc, link.Active, is2025)
		fmt.Printf("\n\n\n\n")
		fmt.Printf("-------------------------------------------------------------\n")
		fmt.Printf("\n\n\n\n")
		if !is2025 && !link.Active {
			continue
		}

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

		// skip already processed session
		if s.Metadata["processed"] == "true" {
			continue
		}

		fmt.Printf("Session ID: %s\n", s.ID)
		customer := new(Customer)
		// Get customer information
		if s.CustomerEmail != "" {
			fmt.Printf("Customer Email: %s\n", s.CustomerEmail)
		}

		if s.CustomerDetails != nil {
			customer.Email = s.CustomerDetails.Email
		}

		if s.Customer != nil {
			customer.Email = s.Customer.Email
		}

		if len(s.CustomFields) > 0 && s.CustomFields[0].Text != nil {
			customer.Name = s.CustomFields[0].Text.Value
		}
		if len(s.CustomFields) > 1 && s.CustomFields[1].Text != nil {
			customer.Company = s.CustomFields[1].Text.Value
		}

		customer.TicketType = linkType(paymentLinkID)
		customer.Force = true

		fmt.Printf("Customer Name: %s -- ", customer.Name)
		fmt.Printf("Customer Company: %s -- ", customer.Company)
		fmt.Printf("Customer Ticket Type: %s\n", customer.TicketType)

		publish(consumerTitle(CUSTOMER), customer)

		// Mark session as processed by adding metadata (optional)
		_, err := session.Update(s.ID, &stripe.CheckoutSessionParams{
			Metadata: map[string]string{
				"processed": "true",
			},
		})
		if err != nil {
			log.Printf("Failed to update session metadata: %v", err)
		}

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
	case "plink_1S4xYEFiHSKRe5D1lDKgFp9Q":
		return "Media pass"
	case "plink_1S5IrSFiHSKRe5D1jALPVk9b":
		return "Stage"
	case "plink_1SDISYFiHSKRe5D1pHeQaX0i":
		return "Exhibitor"
	default:
		return "Standard"
	}
}
