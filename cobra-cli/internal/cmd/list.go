package cmd

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pawelszydlo/humanize"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Lista todos",
	Long:  "Muestra el listado de comandos",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := service.List()
		if err != nil {
			log.Fatal("Error retrieving data ", err)
		}

		h, _ := humanize.New("en")

		tbl := table.NewWriter()

		tbl.SetStyle(table.StyleLight)

		tbl.SetOutputMirror(os.Stdout)

		tbl.AppendHeader(table.Row{"ID", "Name", "Status", "Created At"})

		for _, taskItem := range tasks {
			tbl.AppendRow(table.Row{taskItem.ID, taskItem.Name, taskItem.Status, h.TimeDiffNow(taskItem.CreatedAt, false)})
		}

		tbl.Render()

	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
