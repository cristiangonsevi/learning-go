package user

import "fmt"

type User struct {
	ID   int
	Name string
}

type Repository interface {
	FindOne(id int) (User, error)
	Save(user User) error
	FindAll() ([]User, error)
}

type InMemoryRepository struct {
	data map[int]User
}

func (r *InMemoryRepository) FindOne(id int) (User, error) {
	var user, exist = r.data[id]
	if !exist {
		return User{}, fmt.Errorf("User not found")
	}

	return user, nil
}

func (r *InMemoryRepository) Save(user User) error {
	r.data[user.ID] = user

	return nil
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[int]User),
	}
}
