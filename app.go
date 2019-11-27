package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func registerApp() {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     "https://botsin.space",
		ClientName: "aelita-bot",
		Scopes:     "read write follow",
		Website:    "https://niu.moe/0orpheus",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client-id    : %s\n", app.ClientID)
	fmt.Printf("client-secret: %s\n", app.ClientSecret)
}
