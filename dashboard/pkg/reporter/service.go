package reporter

import (
	"time"
)

const (
	timeout = time.Second * 32
)

type Service interface {
	GetReports(r Report, limit int) []Report
}

type Repository interface {
	FindReports(r Report, limit int) []Report
}

// DBSyncer provides check for update by aggregator
type DBSyncer interface {
	RegForSyncSignal() chan<- chan<- bool
}

type service struct {
	repository Repository
	syncSignal chan<- chan<- bool
}

func NewService(s DBSyncer, r Repository) Service {
	return &service{
		repository: r,
		syncSignal: s.RegForSyncSignal(),
	}
}

func (s *service) GetReports(r Report, limit int) []Report {
	isSynced := make(chan bool, 1) // cap: 1 for non-blocking
	s.syncSignal <- isSynced

	select {
	case <-isSynced:
		return s.repository.FindReports(r, limit)
	case <-time.After(timeout):
		return nil
	}
}
