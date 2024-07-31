package main

import (
	"context"
	"fmt"
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
	user := new(users.User)
	if err := database.Get(ctx, "users/fetch", user, "4b15f797-c7a0-4cb3-8f20-8ec0c3509b0d"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(user, "-----------------")
}
