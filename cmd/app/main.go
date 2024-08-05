package main

import (
	"context"
	"log"
	"socious/src/apps/users"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user []users.User
	if err := database.Get(ctx, "users/fetch", &user, "4b15f797-c7a0-4cb3-8f20-8ec0c3509b0d", "000684be-1cf6-4068-8ea9-ee036ac8cec9"); err != nil {
		log.Fatal(err)
	}

	for _, u := range user {
		if err := database.LoadRelations(&u); err != nil {
			log.Fatal(err)
		}
	}
}
