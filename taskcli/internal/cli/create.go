package cli

import (
	"taskcli/internal/taskcli"
)

func createCommand(name *string) {
	taskName := taskcli.Task{Name: *name}

	taskName.CreateTask()
}
