package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"xii/dashboard/pkg/reporter"
)

const (
	limitMax = 100
	urlBase  = "/api/v1/"
)

func Handler(s reporter.Service) http.Handler {
	router := httprouter.New()

	router.GET(urlBase+"reporter", getReports(s))

	return router
}

func getReports(s reporter.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = limitMax
		}

		reports := s.GetReports(reporter.Report{
			Location: r.URL.Query().Get("location"),
			Device:   r.URL.Query().Get("device"),
		},
			limit,
		)

		w.Header().Set("Content-Type", "application/json")

		e := json.NewEncoder(w)
		e.SetIndent(" ", " ")
		e.Encode(reports)
	}
}
