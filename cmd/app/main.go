package main

import (
	"log"
	"socious/src/apps"
	"socious/src/config"
	"socious/src/database"
	"time"
)

func main() {
	config.Init("config.yml")
	if err := database.Connect(&database.ConnectOption{
		URL:         config.Config.Database,
		SqlDir:      config.Config.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	}); err != nil {
		log.Fatal(err)
	}

	apps.Serve()

}
