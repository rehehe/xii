package reporter

type Service interface {
	GetReports(r Report, limit int) []Report
}

type Repository interface {
	FindReports(r Report, limit int) []Report
}

// DBSynchronizer provides check for update by aggregator
type DBSynchronizer interface {
	SyncSignal() chan<- chan<- bool
}

type service struct {
	repository Repository
	syncSignal chan<- chan<- bool
}

func NewService(s DBSynchronizer, r Repository) Service {
	return &service{
		repository: r,
		syncSignal: s.SyncSignal(),
	}
}

func (s *service) GetReports(r Report, limit int) []Report {
	isSynced := make(chan bool)
	s.syncSignal <- isSynced
	<-isSynced

	return s.repository.FindReports(r, limit)
}
