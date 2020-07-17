package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/RyanSpitz/TwitterDadJokes/src/tweet"
	"github.com/robfig/cron"
)

func main() {
	fmt.Println("Starting Twitter-Dadz-Bot v0.01!")

	creds := tweet.Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("API_KEY"),
		ConsumerSecret:    os.Getenv("API_SECRET_KEY"),
	}

	fmt.Printf("%+v\n", creds)

	client, err := tweet.GetClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	} else {
		fmt.Println("Successfully Connected to Twitter.")
	}

	// Prints the pointer to the client
	fmt.Printf("%+v\n", client)

	uri := os.Getenv("MONGODB_URI")

	// The cron job that tweets out once daily
	c := cron.New()
	c.AddFunc("@daily", func() {
		// Gets joke doc from mongoDB atlas cluster
		j, err := tweet.QueryJokeFromDB(uri)
		if err != nil {
			log.Println("Error getting Query")
			log.Println(err)
		}
		tweet.SendTweet(j, client)
	})

	c.Start()

	// Program will only stop when a kill command is used in the terminal
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	fmt.Println("Got signal:", s)
}
