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

	for _, taskItem := range tasks {
		log.Printf("Tarea %s %s", taskItem.ID, taskItem.Name)
	}
}
