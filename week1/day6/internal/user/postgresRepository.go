package user

type PostgresRepository struct {
	data []User
}

func (r *PostgresRepository) FindOne(id int) (User, error) {
	for _, user := range r.data {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, nil
}

func (r *PostgresRepository) Save(user User) error {
	r.data = append(r.data, user)
	return nil
}

func (r *PostgresRepository) FindAll() ([]User, error) {
	return r.data, nil
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{
		data: []User{},
	}
}
