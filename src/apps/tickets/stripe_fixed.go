package tickets

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/paymentlink"
)

func fetchPaymentLinks() {
	// FIXED: Removed Active filter to get ALL payment links
	params := &stripe.PaymentLinkListParams{}

	// Increase pagination limit for better performance
	params.Limit = stripe.Int64(100)

	// List ALL payment links (active and inactive)
	i := paymentlink.List(params)
	for i.Next() {
		link := i.PaymentLink()
		fmt.Printf("Payment Link ID: %s, URL: %s, Active: %v\n",
			link.ID, link.URL, link.Active)
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

	// Expand to get all necessary data including custom fields
	params.AddExpand("data.customer")
	params.AddExpand("data.customer_details")
	params.AddExpand("data.custom_fields")

	// List checkout sessions for the payment link
	i := session.List(params)
	for i.Next() {
		s := i.CheckoutSession()
		fmt.Printf("Session ID: %s\n", s.ID)
		customer := new(Customer)

		// Initialize with email first
		if s.CustomerEmail != "" {
			customer.Email = s.CustomerEmail
		} else if s.CustomerDetails != nil {
			customer.Email = s.CustomerDetails.Email
		}

		// FIXED: Properly capture from custom fields
		var customFullName string
		var customOrganization string

		// Process all custom fields
		for _, field := range s.CustomFields {
			// Debug logging
			fmt.Printf("Custom Field - Key: %s, Value: %s\n",
				field.Key, field.Text.Value)

			// Check for Full Name field (Custom Field 1)
			// Match by various possible key patterns
			switch field.Key {
			case "fullname", "full_name", "name", "customfield1":
				customFullName = field.Text.Value
			case "organization", "company", "companyname", "org", "customfield2":
				customOrganization = field.Text.Value
			}

			// Also check by label if available
			if field.Label != nil {
				switch field.Label.Custom {
				case "Full Name", "Name", "Customer Name":
					customFullName = field.Text.Value
				case "Organization", "Company", "Company Name":
					customOrganization = field.Text.Value
				}
			}
		}

		// Set customer data with priority to custom fields
		if customFullName != "" {
			customer.Name = customFullName
		} else if s.CustomerDetails != nil && s.CustomerDetails.Name != "" {
			// Fallback to CustomerDetails if custom field not found
			customer.Name = s.CustomerDetails.Name
		} else if s.Customer != nil && s.Customer.Name != "" {
			// Final fallback to Customer object
			customer.Name = s.Customer.Name
		}

		// Set organization/company
		if customOrganization != "" {
			customer.Company = customOrganization
		}

		// Set ticket type
		customer.TicketType = linkType(paymentLinkID)

		// Log the extracted data
		fmt.Printf("Extracted - Name: %s, Email: %s, Company: %s, Type: %s\n",
			customer.Name, customer.Email, customer.Company, customer.TicketType)

		// Only publish if we have minimum required data
		if customer.Email != "" && customer.Name != "" {
			publish(consumerTitle(CUSTOMER), customer)
		} else {
			log.Printf("WARNING: Missing required data for session %s (Name: %s, Email: %s)\n",
				s.ID, customer.Name, customer.Email)
		}
	}

	if err := i.Err(); err != nil {
		log.Printf("Error listing checkout sessions: %v\n", err)
	}
}

// Debug function to identify custom field structure
func debugCustomFields(s *stripe.CheckoutSession) {
	fmt.Println("=== DEBUG: Custom Fields Structure ===")
	fmt.Printf("Session ID: %s\n", s.ID)

	if s.CustomFields == nil || len(s.CustomFields) == 0 {
		fmt.Println("No custom fields found")
		return
	}

	for i, field := range s.CustomFields {
		fmt.Printf("Field %d:\n", i+1)
		fmt.Printf("  Key: %s\n", field.Key)
		if field.Label != nil {
			fmt.Printf("  Label (Custom): %s\n", field.Label.Custom)
			fmt.Printf("  Label (Type): %s\n", field.Label.Type)
		}
		if field.Text != nil {
			fmt.Printf("  Value: %s\n", field.Text.Value)
		}
		if field.Dropdown != nil {
			fmt.Printf("  Dropdown Value: %s\n", field.Dropdown.Value)
		}
		if field.Numeric != nil {
			fmt.Printf("  Numeric Value: %s\n", field.Numeric.Value)
		}
		fmt.Printf("  Type: %s\n", field.Type)
		fmt.Printf("  Optional: %v\n", field.Optional)
		fmt.Println("---")
	}
	fmt.Println("=====================================")
}

// Validation function to count all tickets
func validateTicketCount() {
	totalSessions := 0
	totalActive := 0
	totalInactive := 0

	params := &stripe.PaymentLinkListParams{}
	params.Limit = stripe.Int64(100)

	i := paymentlink.List(params)
	for i.Next() {
		link := i.PaymentLink()
		if link.Active {
			totalActive++
		} else {
			totalInactive++
		}

		sessionParams := &stripe.CheckoutSessionListParams{}
		sessionParams.Filters.AddFilter("payment_link", "", link.ID)
		sessionParams.Filters.AddFilter("status", "", "complete")

		j := session.List(sessionParams)
		for j.Next() {
			totalSessions++
		}
	}

	fmt.Printf("=== Ticket Count Validation ===\n")
	fmt.Printf("Total Payment Links: %d (Active: %d, Inactive: %d)\n",
		totalActive+totalInactive, totalActive, totalInactive)
	fmt.Printf("Total Successful Payments: %d\n", totalSessions)
	fmt.Printf("================================\n")
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
	default:
		return "Standard"
	}
}