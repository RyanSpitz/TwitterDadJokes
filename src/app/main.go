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
	}

	// Prints the pointer to the client
	fmt.Printf("%+v\n", client)

	jokesFile := "internal/jokes.txt"

	fmt.Println("Attepting to read jokes into memory...")
	var m map[int]string
	m = make(map[int]string)

	// Loop over lines in file.
	for index, line := range tweet.ScanByLine(jokesFile) {
		m[index] = line
	}

	fmt.Println("Success")

	// current map key
	count := 0

	c := cron.New()
	c.AddFunc("@daily", func() {
		tweet.SendTweet(m[count], client)

		// Prints the tweet than deletes it from memory
		fmt.Println(m[count])
		delete(m, count)
		count++
	})

	c.Start()

	// Program will only stop when a kill command is used in the terminal
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	fmt.Println("Got signal:", s)
}
