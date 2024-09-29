package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"news_scrapper/cmd/commands"
	"news_scrapper/pkg/config"
	"news_scrapper/pkg/crawler"
	"os"
)

const (
	Title = config.Command
)

func init() {
	viper.SetDefault(config.ConfigurationFile, "./cfg-dev.json")
}

func main() {
	printBanner()
	execCommands()
}

func execCommands() {
	rootCmd := commands.GetRootCmd()
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == config.VersionCmd || cmd.Name() == config.CronCmd || cmd.Name() == config.RemoveCmd || cmd.Name() == "help" {
			return nil
		}
		tmpFile, err := commands.GetConfigFile(cmd)
		if err != nil {
			return err
		}
		settings := config.Settings{}
		err = config.LoadSettingsFile(tmpFile, &settings)
		if err != nil {
			return err
		}
		viper.Set(config.SettingsLoaded, crawler.CreateCrawlerHelper(settings))
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printBanner() {
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(Title)
	fmt.Printf("%s", renderStr)
}
