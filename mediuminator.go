package main

import (
	"flag"
	"fmt"
	"github.com/darkhelmet/twitterstream"
	"log"
	"strings"
)

func main() {
	consumerKey := flag.String("consumer-key", "", "consumer key")
	consumerSecret := flag.String("consumer-secret", "", "consumer secret")
	accessToken := flag.String("access-token", "", "access token")
	accessSecret := flag.String("access-secret", "", "access token secret")
	flag.Parse()

	listenForTweets(*consumerKey, *consumerSecret, *accessToken, *accessSecret)
}

func listenForTweets(consumerKey string, consumerSecret string, accessToken string, accessSecret string) {
	twitter := twitterstream.NewClient(consumerKey, consumerSecret, accessToken, accessSecret)
	conn, err := twitter.Track("medium com")
	if err != nil {
		log.Fatal("unable to connect to twitter", err)
	}
	for {
		if tweet, err := conn.Next(); err == nil {
			for _, url := range tweet.Entities.Urls {
				if !strings.Contains(*url.ExpandedUrl, "medium.com") {
					continue
				}
				tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
				fmt.Println(tweetUrl, tweet.Text, *url.ExpandedUrl)
			}
		}
	}
}
