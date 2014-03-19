package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/edsu/mediator/medium"
	"github.com/eikeon/web"
)

type Resource struct {
	Route *web.Route
	Site  *Site
}

func (r *Resource) Path() string {
	return r.Route.Data["Path"]
}

func (r *Resource) Title() string {
	return r.Route.Data["Title"]
}

func (r *Resource) Type() string {
	return r.Route.Data["Type"]
}

func (r *Resource) Description() string {
	return r.Route.Data["Description"]
}

func (r *Resource) Photo() string {
	return r.Route.Data["Photo"]
}

type Site struct {
	Name    string
	Host_   string `json:"Host"`
	Static  string
	Routes_ []*web.Route `json:"Routes"`
}

func (s *Site) Host() string {
	return s.Host_
}

func (s *Site) Routes() []*web.Route {
	return s.Routes_
}

func (s *Site) GetResource(name string, route *web.Route, vars web.Vars) web.Resource {
	switch name {
	case "home":
		return &Resource{Route: route, Site: s}
	default:
		log.Printf("Warning: unexpected name '%s'\n", route.Name)
		return nil
	}
}

var Address *string
var Root *string

func main() {
	Address := flag.String("address", ":9999", "http service address")
	//Host := flag.String("host", "localhost", "")
	Root = flag.String("root", "dist", "...")
	flag.Parse()

	web.Root = Root

	h := web.NewHub()
	go func() {
		for mention := range medium.Tweets() {
			h.In <- web.Message{"Tweet": mention.Tweet, "Story": mention.Story, "Count": mention.Count}
		}
	}()
	http.Handle("/messages", h.Handler())

	var s Site
	if h, err := web.Handler(&s); err == nil {
		http.Handle("/", h)
		server := &http.Server{Addr: *Address}
		log.Println("starting server on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

}
