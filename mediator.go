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

type Tweet struct {
	Url         string `db:"RANGE"`
	Text        string
	Created     string
	Story       string `db:"HASH"`
	TwitterUser string
}

var db dynamodb.DynamoDB

var tweetTableName string = "mediator-tweet"
var storyTableName string = "mediator-story"

func createTable(name string, i interface{}) {
	db = dynamodb.NewDynamoDB()
	if db != nil {
		t, err := db.Register(name, i)
		if err != nil {
			panic(err)
		}
		pt := dynamodb.ProvisionedThroughput{ReadCapacityUnits: 1, WriteCapacityUnits: 1}
		if _, err := db.CreateTable(t.TableName, t.AttributeDefinitions, t.KeySchema, pt, nil); err != nil {
			log.Println("CreateTable:", err)
		}
		for {
			if description, err := db.DescribeTable(name, nil); err != nil {
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

func init() {
	createTable(tweetTableName, (*Tweet)(nil))
	createTable(storyTableName, (*story.Story)(nil))
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
again:
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

				tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
				created := tweet.CreatedAt.Format(time.RFC3339Nano)
				t := &Tweet{Url: tweetUrl, Text: tweet.Text, Created: created, Story: story.Url, TwitterUser: tweet.User.ScreenName}
				if db != nil {
					db.PutItem(tweetTableName, db.ToItem(t), nil)
				}
				fmt.Printf("%#v\n", t)
				fmt.Printf("%#v\n", story)
				conditions := dynamodb.KeyConditions{"Story": {[]dynamodb.AttributeValue{{"S": t.Story}}, "EQ"}}
				if qr, err := db.Query(tweetTableName, &dynamodb.QueryOptions{KeyConditions: conditions, Select: "COUNT"}); err == nil {
					fmt.Println("number of times story mentioned: ", qr.Count)
				} else {
					log.Println("query error:", err)
				}
				fmt.Println()
			}
		} else {
			log.Println("err: ", err)
			time.Sleep(time.Second)
			goto again
		}
	}
}
