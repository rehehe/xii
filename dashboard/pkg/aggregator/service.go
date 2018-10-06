package aggregator

import (
	"sync"
	"time"
)

var (
	isSyncedInitCap  = 4
	providersInitCap = 2

	rateThreshold = 256
	rateDuration  = time.Millisecond * 250

	timeoutDuration = time.Millisecond * 16000
)

type Service interface {
	Run()
	RegForSyncSignal() chan<- chan<- bool
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

	syncSignalReg chan chan<- bool
	syncSeekers   []chan<- bool
	syncRequested bool
	syncSeekerMut sync.Mutex

	result chan []Report

	rateAcc int

	aggMut sync.Mutex
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
		providers:  make([]Provider, 0, providersInitCap),

		syncSignalReg: make(chan chan<- bool),
		syncSeekers:   make([]chan<- bool, 0, isSyncedInitCap),

		result: make(chan []Report),
	}
}

func (s *service) Run() {
	for {
		select {
		case seeker := <-s.syncSignalReg:
			s.syncSeekerMut.Lock()
			s.syncSeekers = append(s.syncSeekers, seeker)
			s.syncSeekerMut.Unlock()

			if len(s.syncSeekers) == 1 {
				s.syncRequested = true
			}

		case <-time.Tick(rateDuration):
			if s.syncRequested {
				s.syncRequested = false
				s.rateAcc = 0
				s.aggregate()
			} else if rateThreshold < s.rateAcc {
				s.rateAcc = 0
				s.aggregate()
			} else {
				s.rateAcc++
			}
		}
	}
}

func (s *service) aggregate() {
	s.aggMut.Lock()
	defer s.aggMut.Unlock()

	// results buffer to save them to storage at once
	resBuf := make([]Report, 0, 1024)

	// get buffered results from previous round after timeout and save it
	for 0 < len(s.result) {
		reports := <-s.result
		resBuf = append(resBuf, reports...)
	}
	s.setReports(resBuf)

	// signal all workers for update
	providerSignaled := 0
	for _, p := range s.providers {
		select {
		case p.Provide() <- s.getLastResultIDbyProvider(p.Name()):
			providerSignaled++
		default: // non-blocking send
		}
	}

	timeoutTriggered := false
	timeout := time.After(timeoutDuration)
	for 0 < providerSignaled && !timeoutTriggered {
		select {
		case reports := <-s.result:
			resBuf = append(resBuf, reports...)
			providerSignaled--

		case <-timeout:
			timeoutTriggered = true
		}
	}
	s.setReports(resBuf)

	s.syncSeekerMut.Lock()
	for _, synced := range s.syncSeekers {
		synced <- true
	}
	s.syncSeekers = s.syncSeekers[:0]

	s.syncSeekerMut.Unlock()
}

func (s *service) RegisterNewProvider(p Provider) {
	s.providers = append(s.providers, p)
	p.RegisterResultChan(s.result)
	go p.Run()
}

func (s *service) RegForSyncSignal() chan<- chan<- bool {
	return s.syncSignalReg
}

func (s *service) setReports(rs []Report) error {
	if len(rs) != 0 {

	}
	return s.repository.CreateReports(rs)
}

func (s *service) getLastResultIDbyProvider(provider string) uint {
	return s.repository.LastResultIDbyProvider(provider)
}
