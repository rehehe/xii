package aggregator

import (
	"log"
	"time"
)

var (
	isSyncedCap           = 16
	providersCap          = 8
	tickerTimeoutDuration = time.Millisecond * 3000
	tickerRateDuration    = time.Millisecond * 100
)

type Service interface {
	Run()
	SyncSignal() chan<- chan<- bool
	RegisterNewProvider(Provider)
}

type Repository interface {
	CreateReports(rs []Report) error
	LastResultIDbyProvider(provider string) uint
}

type Provider interface {
	Run()
	Stop()
	Name() string
	Provide() chan<- uint
	RegisterResultChan(chan<- []Report)
}

type service struct {
	repository Repository
	providers  []Provider

	syncSignal chan chan<- bool
	isSynced   []chan<- bool

	result chan []Report
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
		providers:  make([]Provider, 0, providersCap),

		syncSignal: make(chan chan<- bool),
		isSynced:   make([]chan<- bool, 0, isSyncedCap),

		result: make(chan []Report),
	}
}

func (s *service) Run() {
	timeout := time.NewTicker(tickerTimeoutDuration)
	rate := time.NewTicker(tickerRateDuration)

	for {
		select {
		case isSynced := <-s.syncSignal:
			s.isSynced = append(s.isSynced, isSynced)

		case reports := <-s.result:
			s.setReports(reports)
			log.Printf("isSynced: %v", reports)
			for k, synced := range s.isSynced {
				if synced != nil {
					synced <- true
					s.isSynced[k] = nil
				}
			}

		case <-rate.C:
			for _, isSync := range s.isSynced {
				if isSync != nil {
					for _, p := range s.providers {
						if p != nil {
							p.Provide() <- s.getLastResultIDbyProvider(p.Name())
						}
					}
				}
			}

		case <-timeout.C:
			for k, isSynced := range s.isSynced {
				if isSynced != nil {
					isSynced <- true
					s.isSynced[k] = nil
				}
			}
		}
	}
}

func (s *service) RegisterNewProvider(src Provider) {
	s.providers = append(s.providers, src)
	src.RegisterResultChan(s.result)
	go src.Run()
}

func (s *service) SyncSignal() chan<- chan<- bool {
	return s.syncSignal
}

func (s *service) setReports(rs []Report) error {
	return s.repository.CreateReports(rs)
}

func (s *service) getLastResultIDbyProvider(provider string) uint {
	return s.repository.LastResultIDbyProvider(provider)
}
