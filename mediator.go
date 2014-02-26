package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/darkhelmet/twitterstream"
	"github.com/edsu/mediator/story"
	"github.com/eikeon/dynamodb"
)

var db dynamodb.DynamoDB

var storyTableName string = "mediator-story"

func init() {
	db = dynamodb.NewDynamoDB()
	if db != nil {
		t, err := db.Register(storyTableName, (*story.Story)(nil))
		if err != nil {
			panic(err)
		}
		pt := dynamodb.ProvisionedThroughput{ReadCapacityUnits: 1, WriteCapacityUnits: 1}
		if _, err := db.CreateTable(t.TableName, t.AttributeDefinitions, t.KeySchema, pt, nil); err != nil {
			log.Println("CreateTable:", err)
		}
		for {
			if description, err := db.DescribeTable(storyTableName, nil); err != nil {
				log.Println("DescribeTable err:", err)
			} else {
				log.Println(description.Table.TableStatus)
				if description.Table.TableStatus == "ACTIVE" {
					break
				}
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Println("WARNING: could not create database to persist stories.")
	}
}

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

				if db != nil {
					db.PutItem(storyTableName, db.ToItem(&story), nil)
				}

				// TODO: create a Tweet here instead of doing this output
				tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
				fmt.Println(tweetUrl, tweet.Text, *url.ExpandedUrl)
				fmt.Printf("%#v\n", story)
				fmt.Println()
			}
		} else {
			log.Println("err: ", err)
			time.Sleep(time.Second)
		}
	}
}
