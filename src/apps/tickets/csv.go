package tickets

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
)

func readCSV(path string) []Customer {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all rows
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var customers []Customer

	// Skip header (assumes first row is header)
	for i, row := range records {
		if i == 0 {
			continue
		}

		if len(row) < 5 {
			log.Printf("Skipping row %d, not enough columns: %v\n", i, row)
			continue
		}

		customer := Customer{
			Name:       strings.TrimSpace(row[22]),
			Email:      strings.TrimSpace(row[9]),
			Company:    strings.TrimSpace(row[24]),
			TicketType: strings.TrimSpace(row[5]),
			Force:      true,
		}
		customers = append(customers, customer)
	}
	return customers
}

func verifyCustomers(customers []Customer) []error {

	// Skip header (assumes first row is header)
	var errs []error
	for i, c := range customers {
		row := i + 2 // +2 to account for header and 0-based index

		if strings.TrimSpace(c.Name) == "" {
			errs = append(errs, fmt.Errorf("name cannot be empty; Row %d: Name is empty", row))
		}

		// Email must not be empty and must be valid
		if strings.TrimSpace(c.Email) == "" {
			errs = append(errs, fmt.Errorf("email cannot be empty; Row %d: Email is empty", row))
		}
		if _, err := mail.ParseAddress(c.Email); err != nil {
			errs = append(errs, fmt.Errorf("invalid email: %w ;Row %d: Invalid email format: %s", err, row, c.Email))
		}

		// TicketType: only letters, no digits, no symbols, not empty
		if strings.TrimSpace(c.TicketType) == "" {
			errs = append(errs, fmt.Errorf("ticket type cannot be empty; Row %d: Ticket type is empty", row))
		}
		// Regex: only letters (a-z, A-Z), no numbers/symbols/spaces
		validTicket := regexp.MustCompile(`^[A-Za-z]+(?:/\s*[A-Za-z]+)*$`)
		if !validTicket.MatchString(c.TicketType) {
			errs = append(errs, fmt.Errorf("ticket type must contain only letters (no spaces, numbers, or symbols); Row %d: Ticket type must contain only letters (no spaces, numbers, or symbols): %s", row, c.TicketType))
		}
	}

	return errs
}
