package commands

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	config2 "news_scrapper/pkg/config"
	"news_scrapper/pkg/crawler"
	"news_scrapper/pkg/style"
)

const Command = config2.Command
const InstallationPath = config2.InstallationPath
const InstallationCompanyPath = InstallationPath + "/" + config2.Company
const APP = InstallationCompanyPath + "/" + Command
const SymLinkAPP = InstallationPath + "/" + Command

func GetRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   Command,
		Short: "Fetch news articles",
		Long:  `Fetch news in bulk from different journal sources`,
	}
	rootCmd.PersistentFlags().String(config2.CfgFileFlag, "", "configuration file path")
	_ = viper.BindPFlag(config2.ConfigurationFile, rootCmd.Flags().Lookup(config2.CfgFileFlag))
	//////////////////////////////////
	rootCmd.AddCommand(GetCronCmd(APP))
	rootCmd.AddCommand(GetUninstallCmd())
	rootCmd.AddCommand(version())
	rootCmd.AddCommand(healthCheck())

	articleCmd := ArticleCmd()
	rootCmd.AddCommand(articleCmd)

	fetchCmd := FetchCmd()
	articleCmd.AddCommand(fetchCmd)

	displayIdCmd := DisplayIdsCmd()
	articleCmd.AddCommand(displayIdCmd)

	optionsCmd := DisplayOptionsCmd()
	articleCmd.AddCommand(optionsCmd)

	return rootCmd
}

func ArticleCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   config2.NewsCmd,
		Short: "Article operations",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}
	return cmd
}

func FetchCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   config2.FetchFlagCmd,
		Short: "Fetch articles from a portal",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings := (viper.Get(config2.SettingsLoaded)).(crawler.CrawlerHelper)
			portalIdSlice, err := cmd.Flags().GetStringSlice(config2.PortalIdFlag)
			if err != nil {
				return err
			} else {
				portalIdSlice = preProcessIds(settings, portalIdSlice)
				for i := range portalIdSlice {
					settings.FetchNews(portalIdSlice[i])
				}
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringSlice(config2.PortalIdFlag, []string{}, "the portal id we want to evaluate")
	_ = cmd.MarkFlagRequired(config2.PortalIdFlag)

	cmd.Flags().Int(config2.ActionPage, 1, "the page to navigate")
	_ = viper.BindPFlag(config2.ActionPage, cmd.Flags().Lookup(config2.ActionPage))

	cmd.Flags().Int(config2.ActionLimitBatchSize, 10, "the number of articles to fetch")
	_ = viper.BindPFlag(config2.ActionLimitBatchSize, cmd.Flags().Lookup(config2.ActionLimitBatchSize))

	cmd.Flags().String(config2.ActionLimitSortByDate, "0", "sort the search by Date")
	_ = viper.BindPFlag(config2.ActionLimitSortByDate, cmd.Flags().Lookup(config2.ActionLimitSortByDate))

	cmd.Flags().String(config2.ActionSectionFlag, "", "Fetch an specific section from the website")
	_ = viper.BindPFlag(config2.ActionSectionFlag, cmd.Flags().Lookup(config2.ActionSectionFlag))

	cmd.Flags().String(config2.ActionQueryValue, "", "Lookup an specific term")
	_ = viper.BindPFlag(config2.ActionQueryValue, cmd.Flags().Lookup(config2.ActionQueryValue))

	cmd.Flags().String(config2.ActionEndIndex, "0", "set the endIndex of the request")
	_ = viper.BindPFlag(config2.ActionEndIndex, cmd.Flags().Lookup(config2.ActionEndIndex))

	cmd.Flags().String(config2.ActionOutSubDir, "", "save the results inside a new directory in the context dir")
	_ = viper.BindPFlag(config2.ActionOutSubDir, cmd.Flags().Lookup(config2.ActionOutSubDir))
	return cmd
}

func preProcessIds(helper crawler.CrawlerHelper, ids []string) []string {
	for _, v := range ids {
		if v == config2.AllFlag {
			return helper.DisplayAllIds()
		}
	}
	return ids
}

func DisplayIdsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   config2.DisplayIdsCmd,
		Short: "Display available ids",
		Run: func(cmd *cobra.Command, args []string) {
			settings := (viper.Get(config2.SettingsLoaded)).(crawler.CrawlerHelper)
			_ = settings.DisplayAllIds()
		},
	}
	return cmd
}

func DisplayOptionsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   config2.DisplayOptionsCmd,
		Short: "Display options",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(viper.GetString(config2.PortalIdFlag)) == 0 {
				return errors.New("you must provide the following flag " + config2.PortalIdFlag)
			}
			settings := (viper.Get(config2.SettingsLoaded)).(crawler.CrawlerHelper)
			settings.DisplayAvailableOptions(viper.GetString(config2.PortalIdFlag))
			return nil
		},
	}
	cmd.Flags().String(config2.PortalIdFlag, "", "set the portal id we want to evaluate")
	_ = cmd.MarkFlagRequired(config2.PortalIdFlag)
	_ = viper.BindPFlag(config2.PortalIdFlag, cmd.Flags().Lookup(config2.PortalIdFlag))

	return cmd
}

func GetConfigFile(cmd *cobra.Command) (string, error) {
	var err error
	cfgFile, _ := cmd.Flags().GetString(config2.CfgFileFlag)
	if len(cfgFile) == 0 {
		cfgFile, err = config2.SaveSettingsFromStdin()
		if err != nil {
			return "", errors.New("you must provide configuration file either by using --config flag or piped config")
		}
	}
	return cfgFile, nil
}

func version() *cobra.Command {
	return &cobra.Command{
		Use:   config2.VersionCmd,
		Short: "show current bin version",
		Run: func(cmd *cobra.Command, args []string) {
			style.PrintF(style.Provisioning, "%v", config2.Version)
		},
	}
}

func healthCheck() *cobra.Command {
	return &cobra.Command{
		Use:   config2.Health,
		Short: "trigger health verification",
		Run: func(cmd *cobra.Command, args []string) {
			style.PrintF(style.Provisioning, "%v", config2.Version)
		},
	}
}
