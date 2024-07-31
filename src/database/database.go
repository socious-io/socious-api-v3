package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sony/gobreaker"
)

var (
	db          *sqlx.DB
	cb          *gobreaker.CircuitBreaker
	cbStateChan chan gobreaker.State
	sqlDir      string
)

type ConnectOption struct {
	URL         string
	SqlDir      string
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
}

// Connect initializes the database connection and sets up the circuit breaker
func Connect(options *ConnectOption) error {
	var err error
	db, err = sqlx.Open("postgres", options.URL)
	if err != nil {
		return fmt.Errorf("could not connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %v", err)
	}
	// Initialize the circuit breaker
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "DBCircuitBreaker",
		MaxRequests: options.MaxRequests,
		Interval:    options.Interval,
		Timeout:     options.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("Circuit breaker state changed from %s to %s\n", from, to)
			if cbStateChan != nil {
				cbStateChan <- to
			}
		},
	})
	sqlDir = options.SqlDir
	return nil
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return db
}

// Close closes the database connection
func Close() {
	if db != nil {
		db.Close()
	}
}
