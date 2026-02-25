package main

import (
	"api-sample/internal/user"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	var repo = user.NewPostgresRepository()
	var service = user.NewService(repo)

	err := service.CreateUser(user.User{
		ID:   1,
		Name: "John Doe",
	})
	if err != nil {
		panic(err)
	}

	err2 := service.CreateUser(user.User{
		ID:   2,
		Name: "Jane Doe",
	})
	if err2 != nil {
		panic(err2)
	}

	users, err := repo.FindAll()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Users: %+v\n", users)

}
