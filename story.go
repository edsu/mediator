package mediuminator

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Story struct {
	Description string
	Title       string
	Url         string
	Author      string
}

func NewStory(mediumUrl string) Story {
	var doc *goquery.Document
	var e error
	var story Story

	if doc, e = goquery.NewDocument(mediumUrl); e != nil {
		log.Fatal(e.Error())
	}

	// use the canonical url on the page
	story.Url, _ = doc.Find("link[rel=\"canonical\"]").Attr("href")
	story.Title = doc.Find("title").Text()
	story.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	story.Author, _ = doc.Find("link[rel=\"author\"]").Attr("href")
	return story
}

func (s *Story) String() string {
	return fmt.Sprintf("<%s> %s by %s -- %s", s.Url, s.Title, s.Author, s.Description)
}
