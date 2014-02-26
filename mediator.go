package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/darkhelmet/twitterstream"
	"github.com/edsu/mediator/story"
)

func main() {
	consumerKey := flag.String("consumer-key", "", "consumer key")
	consumerSecret := flag.String("consumer-secret", "", "consumer secret")
	accessToken := flag.String("access-token", "", "access token")
	accessSecret := flag.String("access-secret", "", "access token secret")
	flag.Parse()

	Tweets(*consumerKey, *consumerSecret, *accessToken, *accessSecret)
}

func Tweets(consumerKey string, consumerSecret string, accessToken string, accessSecret string) {
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

				story := story.NewStory(*url.ExpandedUrl)

				// TODO: create a Tweet here instead of doing this output
				tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
				fmt.Println(tweetUrl, tweet.Text, *url.ExpandedUrl)
				fmt.Printf("%#v\n", story)
				fmt.Println()
			}
		}
	}
}
