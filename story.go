package mediuminator

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Story struct {
	Description string
	Title       string
	Url         string
	Author      string
	ImageUrl    string
}

func NewStory(mediumUrl string) Story {
	var doc *goquery.Document
	var e error
	var story Story

	if doc, e = goquery.NewDocument(mediumUrl); e != nil {
		log.Fatal(e.Error())
	}

	story.Url, _ = doc.Find("link[rel=\"canonical\"]").Attr("href")
	story.Title = doc.Find("title").Text()
	story.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	// TODO: there can be more than one rel="author"
	story.Author, _ = doc.Find("link[rel=\"author\"]").Attr("href")
	story.ImageUrl, _ = doc.Find("meta[property=\"og:image\"]").Attr("content")
	return story
}
