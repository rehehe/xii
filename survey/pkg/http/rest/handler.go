package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"xii/survey/pkg/report"
)

const (
	limitMax = 100
)

func Handler(s report.Service) http.Handler {
	router := httprouter.New()

	router.GET("/api/v1/reporter", getReports(s))
	router.POST("/api/v1/reporter", setReport(s))

	return router
}

func getReports(s report.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = limitMax
		}

		reports := s.GetReports(report.Report{
			Location: r.URL.Query().Get("location"),
			Device:   r.URL.Query().Get("device"),
		},
			limit,
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reports)
	}
}

func setReport(s report.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var newReport report.Report
		err := json.NewDecoder(r.Body).Decode(&newReport)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("json decode error: " + err.Error()))
		}
		defer r.Body.Close()

		err = s.SetReport(newReport)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "not created. err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "created")
	}
}
