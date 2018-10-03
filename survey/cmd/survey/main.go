package main

import (
	"flag"
	"log"
	"net/http"

	"xii/survey/pkg/http/rest"
	"xii/survey/pkg/report"
	"xii/survey/pkg/storage/relational"
)

const (
	appName = "survey"
)

var (
	label    string
	bindAddr string
	dbURI    string
)

func init() {
	flag.StringVar(&label, "label", "", "app label")
	flag.StringVar(&label, "l", "", "app label")

	flag.StringVar(&bindAddr, "bind-address", "localhost:1101", "http bind address")
	flag.StringVar(&bindAddr, "b", "localhost:1101", "bind-address shorthand")

	flag.StringVar(&dbURI, "db-uri", "dbuser:dbpassword@/devdb", "database uri/data-Provider-name")
	flag.StringVar(&dbURI, "d", "dbuser:dbpassword@/devdb", "db-uri shorthand")
	flag.Parse()
}

func dbPrefix() string {
	if label != "" && label != appName {
		return label + "_" + appName
	}
	return appName
}

func main() {
	s := relational.NewStorage(dbURI, dbPrefix())
	defer s.Close()

	r := report.NewService(s)

	log.Printf("web link: http://%s/api/v1/reporter", bindAddr)
	router := rest.Handler(r)
	log.Fatal(http.ListenAndServe(bindAddr, router))
}
