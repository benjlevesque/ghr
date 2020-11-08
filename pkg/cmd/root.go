package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// set by ldflags at build. See magefile
var (
	timestamp  = "<not set>"
	commitHash = "<not set>"
	gitTag     = "<not set>"
)

// Execute is the entry point of the CLI
func Execute() {
	if err := makeRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func makeRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "ghr",
		Short:   "Installs & manages github releases on your machine",
		Version: fmt.Sprintf("%s (%s)", gitTag, timestamp),
	}

	rootCmd.AddCommand(makeCompletionCmd())
	rootCmd.AddCommand(makeInstallCmd())
	rootCmd.AddCommand(makeUpdateCmd())
	rootCmd.AddCommand(makeListCmd())
	rootCmd.AddCommand(makeCheckCmd())

	rootCmd.PersistentFlags().StringP("config", "c", "$HOME/.config/ghr/config.yaml", "Configuration file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().BoolP("version", "v", false, "Display version")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display help")

	return rootCmd
}
