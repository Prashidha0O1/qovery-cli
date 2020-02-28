package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"qovery.go/api"
	"qovery.go/util"
)

var applicationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications",
	Long: `LIST show all available applications within a project and environment. For example:

	qovery application list`,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasFlagChanged(cmd) {
			BranchName = util.CurrentBranchName()
			qoveryYML, err := util.CurrentQoveryYML()
			if err != nil {
				util.PrintError("No qovery configuration file found")
				os.Exit(1)
			}
			ProjectName = qoveryYML.Application.Project
		}

		ShowApplicationList(ProjectName, BranchName)
	},
}

func init() {
	applicationListCmd.PersistentFlags().StringVarP(&ProjectName, "project", "p", "", "Your project name")
	applicationListCmd.PersistentFlags().StringVarP(&BranchName, "branch", "b", "", "Your branch name")

	applicationCmd.AddCommand(applicationListCmd)
}

func ShowApplicationList(projectName string, branchName string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"applications name", "status", "endpoints", "databases", "brokers", "storage"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgMagentaColor, tablewriter.FgWhiteColor})
	applications := api.ListApplications(api.GetProjectByName(projectName).Id, branchName)
	if applications.Results == nil || len(applications.Results) == 0 {
		table.Append([]string{"", "", "", "", "", ""})
	} else {
		for _, a := range applications.Results {
			table.Append([]string{
				a.Name,
				a.Status.CodeMessage,
				a.ConnectionURI,
				intPointerValue(a.TotalDatabases),
				intPointerValue(a.TotalBrokers),
				intPointerValue(a.TotalStorage),
			})
		}
	}

	table.Render()
	fmt.Printf("\n")
}
