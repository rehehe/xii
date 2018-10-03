package report

type Service interface {
	GetReports(r Report, limit int) []Report
	SetReport(r Report) error
}

type Repository interface {
	FindReports(r Report, limit int) []Report
	CreateReport(r Report) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetReports(r Report, limit int) []Report {
	return s.repo.FindReports(r, limit)
}

func (s *service) SetReport(r Report) error {
	return s.repo.CreateReport(r)
}
