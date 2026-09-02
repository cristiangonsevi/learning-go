package services

type AuthServiceInterface interface {
	LoginUser(email string, password string) (string, error)
}

type AuthStruct struct {
	repo AuthServiceInterface
}
