package services

type URLRepositoryInterface interface {
	CreateURL()
}

type URLServiceStruct struct {
	repo URLRepositoryInterface
}

func NewURLService(urlRepo URLRepositoryInterface) *URLServiceStruct {

	return &URLServiceStruct{
		repo: urlRepo,
	}
}

func (s *URLServiceStruct) CreateURL() {}
