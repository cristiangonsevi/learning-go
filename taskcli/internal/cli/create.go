package cli

import (
	"log"
	"taskcli/internal/model"
	"taskcli/internal/task"
)

func createCommand(taskService *task.Service, name *string) {

	if *name == "" {
		log.Fatal("Nombre de la tarea es requerido")
	}

	newTask := model.Task{Name: *name}

	err := taskService.CreateTask(&newTask)

	if err != nil {
		log.Fatal("Error al crear la tarea error: ", err)
	}

	log.Printf("La tarea %s fue creada con exito \n", *name)
}
