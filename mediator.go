package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"code.google.com/p/go.net/websocket"
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

var PERM = regexp.MustCompile("[0-9a-f]{8}~")

type templateData map[string]interface{}

var site *template.Template

func getSite() *template.Template {
	if site == nil {
		site = template.Must(template.ParseFiles(path.Join(*Root, "templates/site.html")))
	}
	return site
}

var Address *string
var Root *string

func main() {
	Address := flag.String("address", ":9999", "http service address")
	Root = flag.String("root", "dist", "...")
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	flag.Parse()

	h := NewHub()
	go h.run()

	go Tweets(consumerKey, consumerSecret, accessToken, accessSecret, h.In)

	static := http.Dir(path.Join(*Root, "static/"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := templateData{}
		if r.URL.Path == "/" {
			d["Found"] = true
		} else {
			upath := r.URL.Path
			if !strings.HasPrefix(upath, "/") {
				upath = "/" + upath
				r.URL.Path = upath
			}
			f, err := static.Open(path.Clean(upath))
			if err == nil {
				defer f.Close()
				d, err1 := f.Stat()
				if err1 != nil {
					http.Error(w, "could not stat file", http.StatusInternalServerError)
					return
				}
				url := r.URL.Path
				if d.IsDir() {
					if url[len(url)-1] != '/' {
						http.Redirect(w, r, url+"/", http.StatusMovedPermanently)
						return
					}
				} else {
					if url[len(url)-1] == '/' {
						http.Redirect(w, r, url[0:len(url)-1], http.StatusMovedPermanently)
						return
					}
				}
				if d.IsDir() {
					w.WriteHeader(http.StatusNotFound)
				} else {
					ttl := int64(0)
					if PERM.MatchString(path.Base(url)) {
						ttl = int64(365 * 86400)
					}
					w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ttl))
					http.ServeContent(w, r, d.Name(), d.ModTime(), f)
					return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}

		d["Path"] = r.URL.Path

		var bw bytes.Buffer
		h := md5.New()
		mw := io.MultiWriter(&bw, h)
		err := getSite().ExecuteTemplate(mw, "html", d)
		if err == nil {
			w.Header().Set("ETag", fmt.Sprintf(`"%x"`, h.Sum(nil)))
			w.Header().Set("Content-Length", fmt.Sprintf("%d", bw.Len()))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(bw.Bytes())
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

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
