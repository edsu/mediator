package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/darkhelmet/twitterstream"
	"github.com/edsu/mediator/medium"
	"github.com/eikeon/dynamodb"
	"github.com/eikeon/web"
)

type Tweet struct {
	Url         string `db:"RANGE"`
	Text        string
	Published   string
	Story       string `db:"HASH"`
	TwitterUser string
	// TODO: record the Medium Collection that was referenced
}

func Tweets(consumerKey string, consumerSecret string, accessToken string, accessSecret string, out chan Message) {
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

				story, err := medium.GetStory(*url.ExpandedUrl)
				if err != nil {
					log.Print(err)
					continue
				}

				mediumUser, err := medium.GetUser(story.Author)
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
					collection, err := medium.GetCollection(story.Collection)
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
					out <- Message{"Story": story, "Tweet": *t, "Count": qr.Count}
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
	createTable(storyTableName, (*medium.Story)(nil))
	createTable(mediumUserTableName, (*medium.User)(nil))
	createTable(collectionTableName, (*medium.Collection)(nil))
}

type Message map[string]interface{}

type Hub struct {
	In   chan Message
	outs []chan Message
	sync.Mutex
}

func NewHub() *Hub {
	h := &Hub{}
	h.In = make(chan Message)
	return h
}
func (h *Hub) run() {
	for m := range h.In {
		for _, out := range h.outs {
			select {
			case out <- m:
			default:
				log.Println("could not broadcast tweet:", m)
			}
		}
	}
}

func (h *Hub) Add(out chan Message) {
	h.Lock()
	h.outs = append(h.outs, out)
	h.Unlock()
}

type Resource struct {
	Route *web.Route
}

func GetResource(route *web.Route, vars web.Vars) web.Resource {
	return &Resource{Route: route}
}

func (r *Resource) Title() string {
	return r.Route.Data["Title"]
}

var Address *string
var Root *string

func main() {
	Address := flag.String("address", ":9999", "http service address")
	Host := flag.String("host", "localhost", "")
	Root = flag.String("root", "dist", "...")
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	flag.Parse()

	h := NewHub()
	go h.run()

	go Tweets(consumerKey, consumerSecret, accessToken, accessSecret, h.In)

	http.Handle("/messages", websocket.Handler(func(ws *websocket.Conn) {
		in := make(chan Message)
		h.Add(in)
		go func() {
			for {
				var m Message
				if err := websocket.JSON.Receive(ws, &m); err == nil {
					//out <- m
				} else {
					log.Println("Message Websocket receive err:", err)
					return
				}
			}
		}()

		for m := range in {
			if err := websocket.JSON.Send(ws, &m); err != nil {
				log.Println("Message Websocket send err:", err)
				break
			}
		}
		ws.Close()
	}))

	web.Root = Root

	getters := web.Getters{
		"home": GetResource,
	}

	if h, err := web.Handler(*Host, getters); err == nil {
		http.Handle("/", h)
		server := &http.Server{Addr: *Address}
		log.Println("starting server on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	addr := *Address
	log.Println("starting:", addr)
	go func(a string) {
		server := &http.Server{Addr: a}
		err := server.ListenAndServe()
		if err != nil {
			log.Print("ListenAndServe:", err)
		}
	}(addr)

	notifyChannel := make(chan os.Signal, 1)
	signal.Notify(notifyChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	sig := <-notifyChannel
	log.Println("handling:", sig)

}
