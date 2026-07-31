package printer

import (
	"cobra-cli/internal/model"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pawelszydlo/humanize"
)

type TaskPrinter struct{}

func New() *TaskPrinter {
	return &TaskPrinter{}
}

func (t *TaskPrinter) PrintTask(tasks []model.TaskModel) {
	h, _ := humanize.New("en")

	tbl := table.NewWriter()

	tbl.SetStyle(table.StyleLight)

	tbl.SetOutputMirror(os.Stdout)

	tbl.AppendHeader(table.Row{"ID", "Name", "Status", "Created At"})

	for _, taskItem := range tasks {
		tbl.AppendRow(table.Row{taskItem.ID, taskItem.Name, taskItem.Status, h.TimeDiffNow(taskItem.CreatedAt, false)})
	}

	tbl.Render()
}
