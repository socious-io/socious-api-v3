package main

import (
	"fmt"
	"log"
	"socious/src/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.Init("config.yml")
	fmt.Println(config.Config)
	m, err := migrate.New(config.Config.MigrationsFile, config.Config.Database)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	log.Println("Migrations applied successfully!")
}
