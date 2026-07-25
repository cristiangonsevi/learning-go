package cli

import (
	"log"
	"taskcli/internal/task"
)

func deleteCommand(taskService *task.Service, taskId *int) {
	if *taskId == 0 {
		log.Fatal("taskid es requerido")
	}

	_, err := taskService.FindTask(taskId)

	if err != nil {
		log.Fatal("Error al borrar la tarea error: ", err)
	}

	taskService.DeleteTask(taskId)

	log.Printf("La tarea %d ha sido borrada correctamente\n", *taskId)
}
