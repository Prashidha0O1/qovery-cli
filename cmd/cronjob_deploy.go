package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/qovery/qovery-client-go"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/qovery/qovery-cli/utils"
)

var cronjobDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a cronjob",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Capture(cmd)

		client := utils.GetQoveryClientPanicInCaseOfError()
		validateCronjobArguments(cronjobName, cronjobNames)
		if cronjobTag != "" && cronjobCommitId != "" {
			utils.PrintlnError(fmt.Errorf("you can't use --tag and --commit-id at the same time"))
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}
		envId := getEnvironmentIdFromContextPanicInCaseOfError(client)

		cronJobList := buildCronJobListFromCronjobNames(client, envId, cronjobName, cronjobNames)
		err := utils.DeployJobs(client, envId, cronJobList, cronjobCommitId, cronjobTag)
		checkError(err)
		utils.Println(fmt.Sprintf("Request to deploy cronjob(s) %s has been queued..", pterm.FgBlue.Sprintf("%s%s", cronjobName, cronjobNames)))
		WatchJobDeployment(client, envId, cronJobList, watchFlag, qovery.STATEENUM_DEPLOYED)
	},
}

func WatchJobDeployment(
	client *qovery.APIClient,
	envId string,
	cronJobs []*qovery.JobResponse,
	watchFlag bool,
	finalServiceState qovery.StateEnum,
) {
	if watchFlag {
		time.Sleep(3 * time.Second) // wait for the deployment request to be processed (prevent from race condition)
		if len(cronJobs) == 1 {
			utils.WatchJob(utils.GetJobId(cronJobs[0]), envId, client)
		} else {
			utils.WatchEnvironment(envId, finalServiceState, client)
		}
	}
}

func init() {
	cronjobCmd.AddCommand(cronjobDeployCmd)
	cronjobDeployCmd.Flags().StringVarP(&organizationName, "organization", "", "", "Organization Name")
	cronjobDeployCmd.Flags().StringVarP(&projectName, "project", "", "", "Project Name")
	cronjobDeployCmd.Flags().StringVarP(&environmentName, "environment", "", "", "Environment Name")
	cronjobDeployCmd.Flags().StringVarP(&cronjobName, "cronjob", "n", "", "Cronjob Name")
	cronjobDeployCmd.Flags().StringVarP(&cronjobNames, "cronjobs", "", "", "Cronjob Names (comma separated) (ex: --cronjobs \"cron1,cron2\")")
	cronjobDeployCmd.Flags().StringVarP(&cronjobCommitId, "commit-id", "c", "", "Cronjob Commit ID")
	cronjobDeployCmd.Flags().StringVarP(&cronjobTag, "tag", "t", "", "Cronjob Tag")
	cronjobDeployCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch cronjob status until it's ready or an error occurs")
}
