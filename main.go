package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Loading config")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/aelitamasto/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetConfigType("yaml")
	host := viper.GetString("host")
	port := viper.GetString("port")
	server := viper.GetString("server")
	clientID := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")
	email := viper.GetString("email")
	password := viper.GetString("password")
	owner := viper.GetString("owner")
	c := mastodon.NewClient(&mastodon.Config{
		Server:       server,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	err = c.Authenticate(context.Background(), email, password)
	if err != nil {
		log.Fatal(err)
	}
	r := regexp.MustCompile(`</a><\/span>(.+)?<\/p>`)
	db := CreateTextDB("/var/db.txt")
	db.Open()
	aelita := connectToAelita(host, port)
	br := true
	for br == true {
		notifications := getDMs(c)
		fmt.Println("Reading notifications")
		for _, n := range notifications {
			if (!db.Check(string(n.ID)) && n.Status.Visibility == (mastodon.VisibilityDirectMessage)) {
				fmt.Println("Responding")
				cmd := parseStatus(n.Status, r)
				err := db.Add(string(n.ID))
				if err == false {
					log.Fatal("Failed to add to db")
				}
				if cmd == "STOP" {
					br = false
				} else {
					fmt.Println("Checking with Aelita")
					respond := sendCommand(aelita,cmd)
					fmt.Println("Sending DM")
					sendDM(c, respond, owner)
					fmt.Println("DM sent")
				}
			}
		}
		if br == true {
			time.Sleep(10 * time.Second)
		}
		fmt.Println("checking again")
	}
	db.Close()
	disconnectFromAelita(aelita)
	fmt.Println("done")
}

func getDMs(c *mastodon.Client) []mastodon.Notification {
	notifications, err := c.GetNotifications(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	var dms []mastodon.Notification
	for i := len(notifications) - 1; i >= 0; i-- {
		if notifications[i].Type == "mention" {
			//status := notifications[i].Status
			//fmt.Println(parseStatus(status,r))
			dms = append(dms, *notifications[i])
		}
	}
	return dms
}

func sendDM(c *mastodon.Client, text string, owner string) {
	content := owner + " " + text
	_, err := c.PostStatus(context.Background(), &mastodon.Toot{
		Status:     content,
		Visibility: mastodon.VisibilityDirectMessage,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func parseStatus(s *mastodon.Status, r *regexp.Regexp) string {
	res := r.FindStringSubmatch((*s).Content)
	return strings.TrimSpace(res[1])
}

type replyDB interface {
	Open() bool
	Close() bool
	Check(string) bool
	Add(string) bool
	Sync() bool
}
