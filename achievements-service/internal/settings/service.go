package settings

type Service interface {
	GetAll() ([]Setting, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAll() ([]Setting, error) {
	return s.repo.GetAll()
}
