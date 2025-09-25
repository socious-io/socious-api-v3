package tickets

import (
	"encoding/csv"
	"log"
	"os"
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
			Name:       strings.TrimSpace(row[1] + " " + row[2]),
			Email:      strings.TrimSpace(row[0]),
			Company:    strings.TrimSpace(row[3]),
			TicketType: strings.TrimSpace(row[4]),
			Force:      true,
		}
		customers = append(customers, customer)
	}
	return customers
}
