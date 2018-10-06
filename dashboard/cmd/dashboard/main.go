package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"xii/dashboard/pkg/aggregator"
	"xii/dashboard/pkg/http/rest"
	"xii/dashboard/pkg/reporter"
	"xii/dashboard/pkg/storage/relational"
	"xii/dashboard/pkg/worker"
)

const (
	appName = "dashboard"
)

var (
	label     string
	bindAddr  string
	dbURI     string
	providers []string
)

func init() {
	flags()
}

func flags() {
	flag.StringVar(&label, "label", "", "app label")
	flag.StringVar(&label, "l", "", "app label")

	flag.StringVar(&bindAddr, "bind-address", "localhost:1100", "http bind address")
	flag.StringVar(&bindAddr, "b", "localhost:1100", "bind-address shorthand")

	flag.StringVar(&dbURI, "db-uri", "dbuser:dbpassword@/devdb", "database uri/data-Provider-name")
	flag.StringVar(&dbURI, "d", "dbuser:dbpassword@/devdb", "db-uri shorthand")

	var providersStr string
	flag.StringVar(&providersStr,
		"providers",
		"blue http://localhost:1101/api/v1/reporter?limit=1000 "+
			"red http://localhost:1102/api/v1/reporter?limit=1000",
		"space separated pair of providersStr/survey-suppliers title and URL")
	flag.StringVar(&providersStr,
		"p",
		"blue http://localhost:1101/api/v1/reporter?limit=1000 "+
			"red http://localhost:1102/api/v1/reporter?limit=1000",
		"providersStr shorthand")

	flag.Parse()

	providers = strings.Split(providersStr, " ")
	if 0 == len(providers) || len(providers)%2 != 0 {
		log.Fatal("providers list is empty or not pair")
	}
}

func dbPrefix() string {
	if label != "" && label != appName {
		return label + "_" + appName
	}
	return appName
}

func main() {
	sto := relational.NewStorage(dbURI, dbPrefix())
	defer sto.Close()

	agg := aggregator.NewService(sto)
	for i := 0; i < len(providers); i += 2 {
		agg.RegisterNewProvider(
			worker.NewService(providers[i], providers[i+1]),
		)
	}
	go agg.Run()

	r := reporter.NewService(agg, sto)

	log.Printf("web link: http://%s/api/v1/reporter",
		strings.Replace(bindAddr, "0.0.0.0", "localhost", 1))
	router := rest.Handler(r)
	log.Fatal(http.ListenAndServe(bindAddr, router))
}
