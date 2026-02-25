package user

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {

	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUser(id int) (User, error) {
	return s.repo.FindOne(id)
}

func (s *Service) CreateUser(user User) error {
	return s.repo.Save(user)
}
