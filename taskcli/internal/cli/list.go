package cli

import (
	"log"
	"taskcli/internal/task"
)

func newListCommand(taskService *task.Service) {
	tasks, err := taskService.ListTasks()

	if err != nil {
		log.Fatal("Error al listar las tareas error: ", err)
	}

	if len(tasks) == 0 {
		log.Println("No hay tareas")
		return
	}

	for _, taskItem := range tasks {
		log.Printf("Tarea %v %s", taskItem.ID, taskItem.Name)
	}
}
