package tweet

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Credentials holds the key info for the twitter API
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Joke is a struct of a mongodb joke doc formated accordingly
type Joke struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Body string             `bson:"body,omitempty"`
	Used bool               `bson:"used"`
}

// GetClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func GetClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

// SendTweet will take in client info and a string
// It will then send a tweet to twitter.
func SendTweet(j Joke, clientTwitter *twitter.Client) {
	tweet, _, err := clientTwitter.Statuses.Update(j.Body, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", tweet)
}

// QueryJokeFromDB will connect to a mongoDB atlas cluster
// then it will query and unused joke and return it
// this joke is then updated in the DB with the used field set to true
func QueryJokeFromDB(URI string) (Joke, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Failed to ping Atlas cluster")
		log.Fatal(err)
	}

	collection := client.Database("Jokes").Collection("Jokes") // Working directory

	var j Joke

	if err = collection.FindOne(ctx, bson.M{"used": false}).Decode(&j); err != nil {
		log.Fatal(err)
	}
	fmt.Println(j)
	j.Used = true
	if err = collection.FindOneAndReplace(ctx, bson.M{"_id": j.ID}, j).Decode(&j); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully writen")
	}

	return j, nil
}
