package worker

import (
	"encoding/json"
	"net/http"

	"xii/dashboard/pkg/aggregator"
)

type Service interface {
	Run()
	Stop()
	Name() string
	Provide() chan<- uint
	RegisterResultChan(chan<- []aggregator.Report)
}

type service struct {
	name   string
	addr   string
	lastID uint

	notify chan uint
	result chan<- []aggregator.Report
	stop   chan bool
}

func NewService(name string, addr string) Service {
	return &service{
		name:   name,
		addr:   addr,
		notify: make(chan uint),
		stop:   make(chan bool),
	}
}

func (s *service) Run() {
	for {
		select {
		case id := <-s.notify:
			if s.lastID < id || s.lastID == 0 {
				rs := s.fetchReport(id)
				if 0 < len(rs) {
					s.lastID = rs[len(rs)-1].ID
					s.result <- rs
				}
			}
		case <-s.stop:
			break
		}
	}
}

func (s *service) Stop() {
	s.stop <- true
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Provide() chan<- uint {
	return s.notify
}

func (s *service) RegisterResultChan(result chan<- []aggregator.Report) {
	s.result = result
}

func (s *service) fetchReport(id uint) []aggregator.Report {
	var rs []aggregator.Report
	resp, err := http.Get(s.addr)
	if err != nil {
		return rs
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&rs)

	for k, v := range rs {
		if id == v.ID {
			rs = rs[k+1:]
			break
		} else if id < v.ID {
			rs = rs[k:]
			break
		}
	}

	for k := range rs {
		rs[k].Provider = s.name
	}

	return rs
}
