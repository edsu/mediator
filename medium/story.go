package medium

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Story struct {
	Description string
	Title       string
	Url         string `db:"HASH"`
	Author      string
	ImageUrl    string
	Published   string
}

func GetStory(mediumUrl string) Story {
	var doc *goquery.Document
	var e error
	var story Story

	if doc, e = goquery.NewDocument(mediumUrl); e != nil {
		log.Fatal(e.Error())
	}

	story.Url, _ = doc.Find("link[rel=\"canonical\"]").Attr("href")
	story.Url = strings.TrimRight(story.Url, "/")

	story.Title = doc.Find("title").Text()
	story.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	// TODO: there can be more than one rel="author", we (minimally) want the medium.com one
	story.Author, _ = doc.Find("link[rel=\"author\"]").Attr("href")
	story.ImageUrl, _ = doc.Find("meta[property=\"og:image\"]").Attr("content")
	story.Published, _ = doc.Find("meta[property=\"article:published_time\"]").Attr("content")

	return story
}
