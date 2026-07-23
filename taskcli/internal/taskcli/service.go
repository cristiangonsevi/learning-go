package taskcli

import "fmt"

func (t *Task) CreateTask() {
	fmt.Println("Creando task service", t.Name)
}
