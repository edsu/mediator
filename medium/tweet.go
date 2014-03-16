package medium

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/darkhelmet/twitterstream"
	"github.com/eikeon/dynamodb"
)

type Tweet struct {
	Url         string `db:"RANGE"`
	Text        string
	Published   string
	Story       string `db:"HASH"`
	TwitterUser string
	// TODO: record the Medium Collection that was referenced
}

type Mention struct {
	Tweet *Tweet
	Story *Story
	Count int
}

func Tweets() <-chan Mention {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	mentions := make(chan Mention)
	go func() {
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

					story, err := GetStory(*url.ExpandedUrl)
					if err != nil {
						log.Print(err)
						continue
					}

					mediumUser, err := GetUser(story.Author)
					if err != nil {
						log.Print(err)
						continue
					}

					if db != nil {
						db.PutItem(storyTableName, db.ToItem(&story), nil)
						db.PutItem(mediumUserTableName, db.ToItem(&mediumUser), nil)
					}

					// not all stories are part of a collection
					if story.Collection != "" {
						collection, err := GetCollection(story.Collection)
						if err == nil && db != nil {
							db.PutItem(collectionTableName, db.ToItem(&collection), nil)
						}
					}

					tweetUrl := "http://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdString
					published := tweet.CreatedAt.Format(time.RFC3339Nano)
					t := &Tweet{Url: tweetUrl, Text: tweet.Text, Published: published, Story: story.Url, TwitterUser: tweet.User.ScreenName}
					if db != nil {
						db.PutItem(tweetTableName, db.ToItem(t), nil)
					}
					fmt.Printf("%#v\n", t)
					fmt.Printf("%#v\n", story)
					conditions := dynamodb.KeyConditions{"Story": {[]dynamodb.AttributeValue{{"S": t.Story}}, "EQ"}}
					if qr, err := db.Query(tweetTableName, &dynamodb.QueryOptions{KeyConditions: conditions, Select: "COUNT"}); err == nil {
						fmt.Println("number of times story mentioned: ", qr.Count)
						mentions <- Mention{Story: &story, Tweet: t, Count: qr.Count}
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
	}()
	return mentions
}

var db dynamodb.DynamoDB

var tweetTableName string = "mediator-tweet"
var storyTableName string = "mediator-story"
var mediumUserTableName string = "mediator-medium-user"
var collectionTableName string = "mediator-collection"

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
	createTable(storyTableName, (*Story)(nil))
	createTable(mediumUserTableName, (*User)(nil))
	createTable(collectionTableName, (*Collection)(nil))
}
