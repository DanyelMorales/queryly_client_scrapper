package commands

import (
	"github.com/spf13/cobra"
	config2 "news_scrapper/pkg/config"
	"news_scrapper/pkg/cron"
	"news_scrapper/pkg/style"
	"os"
)

func GetUninstallCmd() *cobra.Command {
	var uninstall = &cobra.Command{
		Use:   config2.RemoveCmd,
		Short: "crontab operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeBinaries()
		},
	}
	return uninstall
}

func removeBinaries() error {
	style.OnAction("Attempting to remove cron tasks")
	err := cron.UninstallCron("", cron.GetCronExecLine(Placeholder, "*", APP, "", "", ""))
	if err != nil {
		style.ExitActionF("error uninstalling cron: %v", err)
		return err
	}
	style.OnAction("Attempting to remove installation path")
	err = os.RemoveAll(InstallationCompanyPath)
	if err != nil {
		style.FailedActionF("%v", err)
		return err
	}

	style.OnAction("Attempting to remove symlinks")
	err = os.Remove(SymLinkAPP)
	if err != nil {
		style.FailedActionF("%v", err)
		return err
	}

	cfgPath, err := config2.GetLocalConfigPath()
	if err != nil {
		return err
	}
	style.OnAction("Attempting to remove .config directory")
	err = os.RemoveAll(cfgPath)
	if err != nil {
		style.FailedActionF("%v", err)
		return err
	}
	style.SuccessfulAction("uninstalled correctly!")
	os.Exit(0)
	return nil
}
