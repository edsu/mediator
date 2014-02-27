package medium

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Collection struct {
	Url         string
	Title       string
	Description string
	ImageUrl    string
}

func GetCollection(mediumUrl string) Collection {
	var doc *goquery.Document
	var e error
	var coll Collection

	if doc, e = goquery.NewDocument(mediumUrl); e != nil {
		log.Fatal(e.Error())
	}

	coll.Url, _ = doc.Find("link[rel=\"canonical\"]").Attr("href")
	coll.Url = strings.TrimRight(coll.Url, "/")

	coll.Title = doc.Find("title").Text()
	coll.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	coll.ImageUrl, _ = doc.Find("meta[property=\"og:image\"]").Attr("content")

	return coll
}
