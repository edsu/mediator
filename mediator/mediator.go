package main

import (
	"flag"

	"github.com/edsu/mediuminator"
)

func main() {
	consumerKey := flag.String("consumer-key", "", "consumer key")
	consumerSecret := flag.String("consumer-secret", "", "consumer secret")
	accessToken := flag.String("access-token", "", "access token")
	accessSecret := flag.String("access-secret", "", "access token secret")
	flag.Parse()

	mediuminator.Tweets(*consumerKey, *consumerSecret, *accessToken, *accessSecret)
}
