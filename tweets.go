package mediuminator

import (
	"fmt"
	"log"
	"strings"

	"github.com/darkhelmet/twitterstream"
)

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
				tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
				fmt.Println(tweetUrl, tweet.Text, *url.ExpandedUrl)
				story := NewStory(*url.ExpandedUrl)
				fmt.Println(story.String())
			}
		}
	}
}
