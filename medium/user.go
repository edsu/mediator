package medium

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type User struct {
	Url         string `db:"HASH"`
	Name        string
	Description string
	ImageUrl    string
	GoogleUrl   string
	TwitterUrl  string
}

func GetUser(mediumUrl string) User {
	var doc *goquery.Document
	var e error
	var user User

	if doc, e = goquery.NewDocument(mediumUrl); e != nil {
		log.Fatal(e.Error())
	}

	user.Url, _ = doc.Find("link[rel=\"canonical\"]").Attr("href")
	user.Url = strings.TrimRight(user.Url, "/")

	user.Name, _ = doc.Find("meta[name=\"title\"]").Attr("content")
	user.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	user.ImageUrl, _ = doc.Find("meta[property=\"og:image\"]").Attr("content")

	doc.Find("link[rel=\"me\"]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			u, err := url.Parse(href)
			if err == nil {
				if u.Host == "twitter.com" {
					u.Scheme = "https" // force https
					user.TwitterUrl = u.String()
				} else if u.Host == "plus.google.com" {
					user.GoogleUrl = u.String()
				}
			}
		}
	})

	return user
}
