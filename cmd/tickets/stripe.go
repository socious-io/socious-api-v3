package main

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

		// Get customer information
		if s.CustomerEmail != "" {
			fmt.Printf("Customer Email: %s\n", s.CustomerEmail)
		}

		if s.CustomerDetails != nil {
			publish(consumerTitle(CUSTOMER), Customer{
				Name:  s.CustomerDetails.Name,
				Email: s.CustomerDetails.Email,
			})
		}

		if s.Customer != nil {
			publish(consumerTitle(CUSTOMER), Customer{
				Name:  s.Customer.Name,
				Email: s.Customer.Email,
			})
		}

	}

	if err := i.Err(); err != nil {
		log.Printf("Error listing checkout sessions: %v\n", err)
	}
}
