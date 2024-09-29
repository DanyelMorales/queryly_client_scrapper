package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	config2 "news_scrapper/pkg/config"
	"news_scrapper/pkg/cron"
	"news_scrapper/pkg/style"
)

const Placeholder = "__flvctl_syncall__"
const DefaultRule = "1 */15 * * *"

func GetCronCmd(tool string) *cobra.Command {
	var cron = &cobra.Command{
		Use:   config2.CronCmd,
		Short: "crontab operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := config2.GetCurrentSudoUser()
			if err != nil {
				return err
			}
			return cron.ListCron(user, Placeholder)
		},
	}
	cron.PersistentFlags().Bool(config2.DefaultFlag, false, "apply default rule")
	cron.PersistentFlags().String(config2.CfgFileFlag, "", "configuration file path")
	cron.PersistentFlags().String(config2.OptionsFlag, "", "set extra options on command execution")
	cron.PersistentFlags().String(config2.LogRedirectFlag, "", "set absolute path where logs should be saved")
	cron.AddCommand(installCmd(tool))
	cron.AddCommand(uninstallCmd(tool))
	return cron
}

func uninstallCmd(tool string) *cobra.Command {
	var cronCmd = &cobra.Command{
		Use:   config2.RemoveCmd,
		Short: "remove cron rules",
		Long:  "for custom rule just pass it as first arg",
		RunE: func(cmd *cobra.Command, args []string) error {
			cronLine, err := getExecLIneFromCfg(tool, false, cmd, args)
			if err != nil {
				return err
			}
			user, err := config2.GetCurrentSudoUser()
			if err != nil {
				return err
			}
			return cron.UninstallCron(user, cronLine)
		},
	}
	cronCmd.Flags().Bool(config2.AllFlag, false, "remove all rules belonging to flvctl")
	return cronCmd
}

func getExecLIneFromCfg(tool string, useCfg bool, cmd *cobra.Command, args []string) (*cron.Cron, error) {
	all, _ := cmd.Flags().GetBool(config2.AllFlag)
	def, _ := cmd.Flags().GetBool(config2.DefaultFlag)
	var rule string
	if def {
		rule = DefaultRule
	} else if all {
		rule = "*"
	} else if len(args) > 0 {
		rule = args[0]
	} else {
		return nil, errors.New(fmt.Sprintf("flag expected: %s, %s or custom rule as arg0", config2.AllFlag, config2.DefaultFlag))
	}

	style.OnActionF("processing cron rule: %s", rule)
	options, _ := cmd.Flags().GetString(config2.OptionsFlag)
	logRedirect, _ := cmd.Flags().GetString(config2.LogRedirectFlag)

	var configFile string
	if useCfg {
		cfg, err := GetConfigFile(cmd)
		if err != nil {
			return nil, err
		}
		configFile = cfg
	}

	return cron.GetCronExecLine(Placeholder, rule, tool, configFile, logRedirect, options), nil
}

func installCmd(tool string) *cobra.Command {
	var cron = &cobra.Command{
		Use:   config2.InstallCmd + " [cron rule | 'd']",
		Short: "install cron rules",
		Long:  "cron rule is either the first arg provided to this command, otherwise use 'd' char to fallback with " + DefaultRule,
		RunE: func(cmd *cobra.Command, args []string) error {
			cronLine, err := getExecLIneFromCfg(tool, true, cmd, args)
			if err != nil {
				return err
			}
			user, err := config2.GetCurrentSudoUser()
			if err != nil {
				return err
			}
			isInit, err := cron.IsInit(user)
			if err != nil {
				return err
			}
			if !isInit {
				err := cron.InitCron(user)
				if err != nil {
					return err
				}
			}
			return cron.InstallCron(user, cronLine)
		},
	}
	return cron
}
