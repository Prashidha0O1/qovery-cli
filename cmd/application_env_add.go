package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"qovery.go/api"
	"qovery.go/util"
)

var applicationEnvAddCmd = &cobra.Command{
	Use:   "add <key> <value>",
	Short: "Add environment variable",
	Long: `ADD an environment variable to an application. For example:

	qovery application env add`,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasFlagChanged(cmd) {
			qoveryYML, err := util.CurrentQoveryYML()
			if err != nil {
				util.PrintError("No qovery configuration file found")
				os.Exit(1)
			}
			BranchName = util.CurrentBranchName()
			ProjectName = qoveryYML.Application.Project
		}

		if len(args) != 2 {
			_ = cmd.Help()
			return
		}

		projectId := api.GetProjectByName(ProjectName).Id
		repositoryId := api.GetRepositoryByCurrentRemoteURL(projectId).Id
		environment := api.GetEnvironmentByBranchId(projectId, repositoryId, BranchName)
		api.CreateApplicationEnvironmentVariable(api.EnvironmentVariable{Key: args[0], Value: args[1]}, projectId, repositoryId,
			environment.Id, environment.Application.Id)

		fmt.Println("ok")
	},
}

func init() {
	applicationEnvAddCmd.PersistentFlags().StringVarP(&ProjectName, "project", "p", "", "Your project name")
	applicationEnvAddCmd.PersistentFlags().StringVarP(&BranchName, "branch", "b", "", "Your branch name")
	// TODO select application

	applicationEnvCmd.AddCommand(applicationEnvAddCmd)
}
